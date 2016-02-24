package server

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"net"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/exago/svc/config"
	"github.com/gorilla/context"
	"github.com/julienschmidt/httprouter"
)

// Router wraps httprouter.Router.
type Router struct {
	*httprouter.Router
}

// Get wraps httprouter's GET.
func (r *Router) Get(path string, handler http.Handler) {
	r.GET(path, wrapHandler(handler))
}

// NewRouter returns a wrapped httprouter.
func NewRouter() *Router {
	r := &Router{httprouter.New()}
	r.NotFound = cstNotFound{}
	r.MethodNotAllowed = cstMethodNotAllowed{}
	return r
}

func wrapHandler(h http.Handler) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		context.Set(r, "params", ps)
		h.ServeHTTP(w, r)
	}
}

type cstMethodNotAllowed struct{}
type cstNotFound struct{}

// ServeHTTP is a custom handler for the "method not allowed" error
func (wrt cstMethodNotAllowed) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	writeData(w, r, http.StatusMethodNotAllowed, "Method not allowed")
}

// ServeHTTP is a custom handler for the 404 error
func (wrt cstNotFound) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	writeData(w, r, http.StatusNotFound, "Not found")
}

// badRequest is handled by setting the status code in the reply to StatusBadRequest.
type badRequest struct{ error }

// notFound is handled by setting the status code in the reply to StatusNotFound.
type notFound struct{ error }

// send is a shorthand
func send(w http.ResponseWriter, r *http.Request, data interface{}, err error) {
	if err == nil {
		writeData(w, r, http.StatusOK, data)
	} else {
		writeError(w, r, err)
	}
}

// outputJSON handles the response for each endpoint.
// It follows the JSEND standard for JSON response.
// See https://labs.omniti.com/labs/jsend
func writeData(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	var err error

	w.Header().Set("Access-Control-Allow-Origin", config.Values.AllowOrigin)
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
		writeError(w, r, err)
		return
	}

	if err = gz.Close(); err != nil {
		writeError(w, r, err)
		return
	}

	if _, err = w.Write(b.Bytes()); err != nil {
		writeError(w, r, err)
	}
}

func writeError(w http.ResponseWriter, r *http.Request, err error) {
	logger := context.Get(r, "logger").(*log.Entry)
	logger.Error(err)

	switch err.(type) {
	case badRequest:
		writeData(w, r, http.StatusBadRequest, err.Error())
	case notFound:
		writeData(w, r, http.StatusNotFound, err.Error())
	default:
		writeData(w, r, http.StatusInternalServerError, err.Error())
	}
}

func getIP(s string) string {
	if ip, _, err := net.SplitHostPort(s); err == nil {
		return ip
	}
	if ip := net.ParseIP(s); ip != nil {
		return ip.To4().String()
	}
	log.Error("Couldn't parse IP %s", s)
	return ""
}
