package main

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/context"
	"github.com/jgautheron/exago-service/badge"
	"github.com/jgautheron/exago-service/config"
	"github.com/jgautheron/exago-service/datafetcher"
	"github.com/jgautheron/exago-service/logger"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

var (
	errInvalidParameter  = errors.New("Invalid parameter passed")
	errInvalidRepository = errors.New("The repository doesn't contain Go code")
	errBadgeError        = errors.New("The badge cannot be fetched")

	// Command line flags
	flagVersion = flag.Bool("version", false, "show the version")

	// Build vars
	// Do not set these manually! these variables
	// are meant to be set through ldflags
	buildTag  string
	buildDate string
)

func init() {
	var err error

	err = config.SetUp()
	if err != nil {
		log.Fatal(err)
	}

	err = logger.SetUp()
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	flag.Parse()
	if *flagVersion {
		fmt.Fprintf(os.Stdout, "Running version "+buildTag+" built on "+buildDate)
		return
	}

	repoHandlers := alice.New(context.ClearHandler, recoverHandler, setLogger, checkRegistry, checkValidRepository, checkCache)
	router := NewRouter()

	router.Get("/:registry/:username/:repository/badge/:type", repoHandlers.ThenFunc(badgeHandler))
	router.Get("/:registry/:username/:repository/valid", repoHandlers.ThenFunc(repoValidHandler))
	router.Get("/:registry/:username/:repository/loc", repoHandlers.ThenFunc(lambdaHandler))
	router.Get("/:registry/:username/:repository/imports", repoHandlers.ThenFunc(lambdaHandler))
	router.Get("/:registry/:username/:repository/lint/:linter", repoHandlers.ThenFunc(lambdaHandler))
	router.Get("/:registry/:username/:repository/test", repoHandlers.ThenFunc(testHandler))
	// router.Get("/:registry/:username/:repository/rank", repoHandlers.ThenFunc(rankHandler))
	router.Get("/:registry/:username/:repository/contents/*path", repoHandlers.ThenFunc(fileHandler))

	hp := config.Get("HttpPort")
	log.Info("Listening on port " + hp)
	if err := http.ListenAndServe(":"+hp, router); err != nil {
		log.Fatal(err)
	}
}

func badgeHandler(w http.ResponseWriter, r *http.Request) {
	ps := context.Get(r, "params").(httprouter.Params)
	tp := ps.ByName("type")

	rp := fmt.Sprintf("%s/%s/%s", ps.ByName("registry"), ps.ByName("username"), ps.ByName("repository"))
	switch tp {
	case "imports":
		out, err := datafetcher.GetImports(rp)
		if err != nil {
			badge.WriteError(w, tp)
			return
		}
		var data []string
		err = json.Unmarshal(*out, &data)
		if err != nil {
			badge.WriteError(w, tp)
			return
		}

		ln := strconv.Itoa(len(data))
		badge.Write(w, "imports", ln, "blue")
	}
}

func rankHandler(w http.ResponseWriter, r *http.Request) {
	ps := context.Get(r, "params").(httprouter.Params)
	log.Info(ps)
}

func lambdaHandler(w http.ResponseWriter, r *http.Request) {
	ps := context.Get(r, "params").(httprouter.Params)

	lambdaFn := strings.Split(r.URL.String()[1:], "/")[3]
	rp := fmt.Sprintf("%s/%s/%s", ps.ByName("registry"), ps.ByName("username"), ps.ByName("repository"))

	var out *json.RawMessage
	var err error

	switch lambdaFn {
	case "imports":
		out, err = datafetcher.GetImports(rp)
	case "loc":
		out, err = datafetcher.GetCodeStats(rp)
	case "lint":
		out, err = datafetcher.GetLintResults(rp, ps.ByName("linter"))
	}

	send(w, r, out, err)
}

func repoValidHandler(w http.ResponseWriter, r *http.Request) {
	ps := context.Get(r, "params").(httprouter.Params)
	username, repository := ps.ByName("username"), ps.ByName("repository")
	code, err := checkRepo(r, username, repository)
	outputJSON(w, r, code, err)
}

func fileHandler(w http.ResponseWriter, r *http.Request) {
	ps := context.Get(r, "params").(httprouter.Params)
	lgr := context.Get(r, "logger").(*log.Entry)

	rp := fmt.Sprintf("%s/%s/%s", ps.ByName("registry"), ps.ByName("username"), ps.ByName("repository"))
	out, err := datafetcher.GetFileContents(lgr, rp, ps.ByName("path")[1:])
	send(w, r, string(out), err)
}

func testHandler(w http.ResponseWriter, r *http.Request) {
	ps := context.Get(r, "params").(httprouter.Params)
	lgr := context.Get(r, "logger").(*log.Entry)

	rp := fmt.Sprintf("%s/%s/%s", ps.ByName("registry"), ps.ByName("username"), ps.ByName("repository"))
	out, err := datafetcher.GetCoverage(lgr, rp)
	send(w, r, out, err)
}

// checkRepo ensures that the matching username/repository exist on Github.
func checkRepo(r *http.Request, username, repository string) (int, error) {
	var err error
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/languages?access_token=%s", username, repository, config.Get("GithubAccessToken"))

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
		lgr := context.Get(r, "logger").(*log.Entry)
		lgr.Info("Not a Go repository")
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
	const timeout = 2 * time.Hour

	ps := context.Get(r, "params").(httprouter.Params)
	lgr := context.Get(r, "logger").(*log.Entry)

	idfr := getCacheIdentifier(r)
	k := fmt.Sprintf("%s/%s/%s", ps.ByName("registry"), ps.ByName("username"), ps.ByName("repository"))
	if _, err := pool.Get().Do("HMSET", k, idfr, output); err != nil {
		lgr.Error(err)
		return
	}
	if _, err := pool.Get().Do("EXPIRE", k, timeout); err != nil {
		lgr.Error(err)
		return
	}
}

func send(w http.ResponseWriter, r *http.Request, data interface{}, err error) {
	if err == nil {
		outputJSON(w, r, http.StatusOK, data)
	} else {
		handleError(w, r, err)
	}
}

// outputJSON handles the response for each endpoint.
// It follows the JSEND standard for JSON response.
// See https://labs.omniti.com/labs/jsend
func outputJSON(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	var err error

	// w.Header().Set("Access-Control-Allow-Origin", config.Get("AllowOrigin"))
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
		// cacheDynamo(r, b.Bytes())
	}

	if _, err = w.Write(b.Bytes()); err != nil {
		handleError(w, r, err)
	}
}
