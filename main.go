package main

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/context"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

const (
	reAlphaNumeric = `^[\w\d_]+$`
)

var (
	invalidParameter = errors.New("Invalid parameter passed")
	lintFail         = errors.New("The linting process didn't finish")
	cfg              *config
)

type config struct {
	githubAccessToken, awsAccessKeyID, awsSecretAccessKey string
}

func init() {
	cfg = &config{
		// Should be overridable later by a logged in user
		githubAccessToken:  os.Getenv("GITHUB_ACCESS_TOKEN"),
		awsAccessKeyID:     os.Getenv("AWS_ACCESS_KEY_ID"),
		awsSecretAccessKey: os.Getenv("AWS_SECRET_ACCESS_KEY"),
	}
}

func main() {
	repoHandlers := alice.New(context.ClearHandler, recoverHandler, setLogger, checkRegistry, checkValidRepository)
	router := NewRouter()

	router.Get("/:registry/:username/:repository/exists", repoHandlers.ThenFunc(repoExistsHandler))
	router.Get("/:registry/:username/:repository/loc", repoHandlers.ThenFunc(lambdaHandler))
	router.Get("/:registry/:username/:repository/imports", repoHandlers.ThenFunc(lambdaHandler))
	router.Get("/:registry/:username/:repository/lint/:linter", repoHandlers.ThenFunc(lambdaHandler))
	router.Get("/:registry/:username/:repository/test", repoHandlers.ThenFunc(testHandler))
	router.Get("/:registry/:username/:repository/contents/*path", repoHandlers.ThenFunc(fileHandler))

	http.ListenAndServe(":8080", router)

	// TODO:
	// - check repository exists caching
	// - caching support
}

func lambdaHandler(w http.ResponseWriter, r *http.Request) {
	ps := context.Get(r, "params").(httprouter.Params)
	lambdaFn := strings.Split(r.URL.String()[1:], "/")[3]

	ctxt := lambdaContext{
		Registry:   ps.ByName("registry"),
		Username:   ps.ByName("username"),
		Repository: ps.ByName("repository"),
	}

	// Specific handling depending of the lambda function
	switch lambdaFn {
	case "lint":
		ctxt.Linters = ps.ByName("linter")
	}

	resp, err := callLambdaFn(lambdaFn, ctxt)
	if err != nil {
		handleError(w, r, err)
		return
	}

	logger := context.Get(r, "logger").(*log.Entry)
	logger.Infoln(lambdaFn, resp.Metadata)

	handleOutput(w, http.StatusOK, resp.Data)
}

func repoExistsHandler(w http.ResponseWriter, r *http.Request) {
	code := http.StatusOK

	ps := context.Get(r, "params").(httprouter.Params)
	username, repository := ps.ByName("username"), ps.ByName("repository")

	// Check if the specified repository is valid
	if err := checkRepo(r, username, repository); err != nil {
		code = http.StatusNotFound
	}

	handleOutput(w, code, nil)
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
	defer resp.Body.Close()

	raw, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		handleError(w, r, err)
		return
	}

	var data map[string]interface{}
	err = json.Unmarshal(raw, &data)
	if err != nil {
		handleError(w, r, err)
		return
	}

	contents, err := base64.StdEncoding.DecodeString(data["content"].(string))
	if err != nil {
		handleError(w, r, err)
		return
	}

	handleOutput(w, http.StatusOK, string(contents))
}

func testHandler(w http.ResponseWriter, r *http.Request) {
	ps := context.Get(r, "params").(httprouter.Params)
	registry, username, repository := ps.ByName("registry"), ps.ByName("username"), ps.ByName("repository")

	fp := fmt.Sprintf("%s/%s/%s", registry, username, repository)
	out, err := exec.Command("docker", "run", "--rm", "-a", "stdout", "-a", "stderr", "exago-runner", fp).CombinedOutput()
	if err != nil {
		handleError(w, r, err)
		return
	}

	resp := map[string]interface{}{}
	err = json.Unmarshal(out, &resp)

	handleOutput(w, http.StatusOK, resp)
}

// checkRepo ensures that the matching username/repository exist on Github.
func checkRepo(r *http.Request, username, repository string) error {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s?access_token=%s", username, repository, cfg.githubAccessToken)

	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("Repository %s not found in Github.", repository)
	}

	return nil
}

// handleOutput handles the response for each endpoint.
// It follows the JSEND standard for JSON response.
// See https://labs.omniti.com/labs/jsend
func handleOutput(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	success := false
	if code == 200 {
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

	json.NewEncoder(w).Encode(res)
}
