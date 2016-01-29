package main

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/fsouza/go-dockerclient"
	"github.com/gorilla/context"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

var (
	errInvalidParameter  = errors.New("Invalid parameter passed")
	errInvalidRepository = errors.New("The repository doesn't contain Go code")
	errEmptyOutput       = errors.New("The container output is empty")

	cfg *config

	// Command line flags
	flagVersion = flag.Bool("version", false, "show the version")

	// Build vars
	// Do not set these manually! these variables
	// are meant to be set through ldflags
	buildTag  string
	buildDate string
)

type config struct {
	githubAccessToken, awsAccessKeyID, awsSecretAccessKey string
	runnerImageName                                       string
	httpPort, redisHost                                   string
}

func init() {
	cfg = &config{
		// Should be overridable later by a logged in user
		githubAccessToken:  os.Getenv("GITHUB_ACCESS_TOKEN"),
		awsAccessKeyID:     os.Getenv("AWS_ACCESS_KEY_ID"),
		awsSecretAccessKey: os.Getenv("AWS_SECRET_ACCESS_KEY"),
		runnerImageName:    os.Getenv("RUNNER_IMAGE_NAME"),
		httpPort:           os.Getenv("HTTP_PORT"),
		redisHost:          os.Getenv("REDIS_HOST"),
	}
}

func main() {
	flag.Parse()
	if *flagVersion {
		fmt.Fprintf(os.Stdout, "Running version "+buildTag+" built on "+buildDate)
		return
	}

	// Basic validation
	if cfg.githubAccessToken == "" ||
		cfg.awsAccessKeyID == "" ||
		cfg.awsSecretAccessKey == "" ||
		cfg.httpPort == "" ||
		cfg.redisHost == "" {
		log.Fatal("Missing environment variable(s)")
	}

	repoHandlers := alice.New(context.ClearHandler, recoverHandler, setLogger, checkRegistry, checkValidRepository, checkCache)
	router := NewRouter()

	router.Get("/:registry/:username/:repository/valid", repoHandlers.ThenFunc(repoValidHandler))
	router.Get("/:registry/:username/:repository/loc", repoHandlers.ThenFunc(lambdaHandler))
	router.Get("/:registry/:username/:repository/imports", repoHandlers.ThenFunc(lambdaHandler))
	router.Get("/:registry/:username/:repository/lint/:linter", repoHandlers.ThenFunc(lambdaHandler))
	router.Get("/:registry/:username/:repository/test", repoHandlers.ThenFunc(testHandler))
	router.Get("/:registry/:username/:repository/contents/*path", repoHandlers.ThenFunc(fileHandler))

	log.Info("Listening on port " + cfg.httpPort)
	if err := http.ListenAndServe(":"+cfg.httpPort, router); err != nil {
		log.Fatal(err)
	}
}

func lambdaHandler(w http.ResponseWriter, r *http.Request) {
	ps := context.Get(r, "params").(httprouter.Params)
	lambdaFn := strings.Split(r.URL.String()[1:], "/")[3]
	linter := ps.ByName("linter")

	ctxt := lambdaContext{
		Registry:   ps.ByName("registry"),
		Username:   ps.ByName("username"),
		Repository: ps.ByName("repository"),
	}

	// Specific handling depending of the lambda function
	switch lambdaFn {
	case "lint":
		ctxt.Linters = linter
	}

	resp, err := callLambdaFn(lambdaFn, ctxt)
	if err != nil {
		handleError(w, r, err)
		return
	}

	logger := context.Get(r, "logger").(*log.Entry)
	logger.Infoln(lambdaFn, resp.Metadata)

	outputJSON(w, r, http.StatusOK, resp.Data)
}

func repoValidHandler(w http.ResponseWriter, r *http.Request) {
	ps := context.Get(r, "params").(httprouter.Params)
	username, repository := ps.ByName("username"), ps.ByName("repository")
	code, err := checkRepo(r, username, repository)
	outputJSON(w, r, code, err)
}

func fileHandler(w http.ResponseWriter, r *http.Request) {
	ps := context.Get(r, "params").(httprouter.Params)
	username, repository, path := ps.ByName("username"), ps.ByName("repository"), ps.ByName("path")

	// Only github implemented for now
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/contents/%s?access_token=%s", username, repository, path, cfg.githubAccessToken)

	resp, err := http.Get(url)
	if err != nil {
		handleError(w, r, err)
		return
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			log.Error(err)
		}
	}()

	raw, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		handleError(w, r, err)
		return
	}

	var data map[string]interface{}
	if err = json.Unmarshal(raw, &data); err != nil {
		handleError(w, r, err)
		return
	}

	contents, err := base64.StdEncoding.DecodeString(data["content"].(string))
	if err != nil {
		handleError(w, r, err)
		return
	}

	outputJSON(w, r, http.StatusOK, string(contents))
}

func testHandler(w http.ResponseWriter, r *http.Request) {
	ps := context.Get(r, "params").(httprouter.Params)
	logger := context.Get(r, "logger").(*log.Entry)
	registry, username, repository := ps.ByName("registry"), ps.ByName("username"), ps.ByName("repository")

	resp := make(chan interface{})
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()

		client, _ := docker.NewClientFromEnv()
		config := docker.Config{
			Image: cfg.runnerImageName,
			Cmd: []string{
				fmt.Sprintf("%s/%s/%s", registry, username, repository),
			},
		}

		rp := fmt.Sprintf("%s-%s-%s", registry, username, repository)
		opts := docker.CreateContainerOptions{Name: "exago-runner-" + rp, Config: &config}
		container, err := client.CreateContainer(opts)
		if err != nil {
			log.Errorln(err, "exago-runner-"+rp)
			resp <- err
			return
		}

		defer func() {
			err = client.RemoveContainer(docker.RemoveContainerOptions{ID: container.ID})
			if err != nil {
				logger.Error(err)
				return
			}
		}()

		err = client.StartContainer(container.ID, &docker.HostConfig{})
		if err != nil {
			resp <- err
			return
		}

		var buf bytes.Buffer
		err = client.Logs(docker.LogsOptions{
			Container:    container.ID,
			OutputStream: &buf,
			Follow:       true,
			Stdout:       true,
			Stderr:       true,
			Timestamps:   false,
		})

		if buf.Len() == 0 {
			resp <- errEmptyOutput
			return
		}

		var msg *json.RawMessage
		if err = json.Unmarshal(buf.Bytes(), &msg); err != nil {
			resp <- err
			return
		}

		log.Info("output sent")

		resp <- msg
	}()

	out := <-resp
	switch out.(type) {
	case error:
		handleError(w, r, out.(error))
	case *json.RawMessage:
		outputJSON(w, r, http.StatusOK, out)
	}
	wg.Wait()
}

// checkRepo ensures that the matching username/repository exist on Github.
func checkRepo(r *http.Request, username, repository string) (int, error) {
	var err error
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/languages?access_token=%s", username, repository, cfg.githubAccessToken)

	resp, err := http.Get(url)
	if err != nil {
		return http.StatusServiceUnavailable, err
	}

	if resp.StatusCode != 200 {
		return http.StatusNotFound, fmt.Errorf("Repository %s not found in Github.", repository)
	}

	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return http.StatusNoContent, err
	}

	data := map[string]int{}
	err = json.Unmarshal(contents, &data)
	if err != nil {
		return http.StatusNoContent, err
	}

	if _, found := data["Go"]; !found {
		logger := context.Get(r, "logger").(*log.Entry)
		logger.Info("Not a Go repository")
		return http.StatusNotAcceptable, errInvalidRepository
	}

	return http.StatusOK, nil
}

func getCacheIdentifier(r *http.Request) string {
	ps := context.Get(r, "params").(httprouter.Params)
	sp := strings.Split(r.URL.String(), "/")
	switch action := sp[4]; {
	case ps.ByName("linter") != "":
		return action + ":" + ps.ByName("linter")
	case ps.ByName("path") != "":
		return action + ":" + ps.ByName("path")
	default:
		return action
	}
}

func cacheOutput(r *http.Request, output []byte) {
	const timeout = 60 * time.Minute

	ps := context.Get(r, "params").(httprouter.Params)
	logger := context.Get(r, "logger").(*log.Entry)

	idfr := getCacheIdentifier(r)
	k := fmt.Sprintf("%s/%s/%s", ps.ByName("registry"), ps.ByName("username"), ps.ByName("repository"))
	if _, err := pool.Get().Do("HMSET", k, idfr, output); err != nil {
		logger.Error(err)
		return
	}
	if _, err := pool.Get().Do("EXPIRE", k, timeout); err != nil {
		logger.Error(err)
		return
	}
}

// outputJSON handles the response for each endpoint.
// It follows the JSEND standard for JSON response.
// See https://labs.omniti.com/labs/jsend
func outputJSON(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	var err error

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Encoding", "gzip")
	w.WriteHeader(code)

	success := false
	if code == http.StatusOK {
		success = true
	}

	// JSend has three possible statuses: success, fail and error
	// In case of error, there is no data sent, only an error message.
	status := "success"
	dataType := "data"
	if !success {
		status = "error"
		dataType = "message"
	}

	res := map[string]interface{}{"status": status}
	if data != nil {
		res[dataType] = data
	}

	// Assuming here that all browsers support gzip encoding
	var b bytes.Buffer
	gz := gzip.NewWriter(&b)
	if err = json.NewEncoder(gz).Encode(res); err != nil {
		handleError(w, r, err)
		return
	}

	if err = gz.Close(); err != nil {
		handleError(w, r, err)
		return
	}

	if success {
		cacheOutput(r, b.Bytes())
	}

	if _, err = w.Write(b.Bytes()); err != nil {
		handleError(w, r, err)
	}
}

func outputFromCache(w http.ResponseWriter, r *http.Request, code int, output []byte) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Encoding", "gzip")
	w.WriteHeader(code)

	if _, err := w.Write(output); err != nil {
		handleError(w, r, err)
	}
}
