package main

import (
	"net/http"

	"github.com/gorilla/context"
	"github.com/julienschmidt/httprouter"
)

// - handler return data, err (output automatically jsend)
// - jsend
// - context

type router struct {
	*httprouter.Router
}

func (r *router) Get(path string, handler http.Handler) {
	r.GET(path, wrapHandler(handler))
}

func NewRouter() *router {
	r := &router{httprouter.New()}
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
	outputJSON(w, r, http.StatusMethodNotAllowed, "Method not allowed")
}

// ServeHTTP is a custom handler for the 404 error
func (wrt cstNotFound) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	outputJSON(w, r, http.StatusNotFound, "Not found")
}
