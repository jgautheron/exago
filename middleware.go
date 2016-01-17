package main

import (
	"fmt"
	"net"
	"net/http"
	"regexp"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/context"
	"github.com/julienschmidt/httprouter"
)

func recoverHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("panic: %+v", err)
				http.Error(w, http.StatusText(500), 500)
			}
		}()

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func checkRegistry(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		params := context.Get(r, "params").(httprouter.Params)
		registry := params.ByName("registry")
		if registry != "github.com" {
			outputJSON(w, r, http.StatusNotImplemented, fmt.Sprintf("the registry %s is not yet supported", registry))
			return
		}
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func checkValidRepository(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		params := context.Get(r, "params").(httprouter.Params)
		username, repository := params.ByName("username"), params.ByName("repository")
		re := regexp.MustCompile(reAlphaNumeric)
		if !re.MatchString(username) || !re.MatchString(repository) {
			handleError(w, r, errInvalidParameter)
			return
		}
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func checkCache(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ps := context.Get(r, "params").(httprouter.Params)
		lgr := context.Get(r, "logger").(*log.Entry)

		idfr := getCacheIdentifier(r)
		k := fmt.Sprintf("%s/%s/%s", ps.ByName("registry"), ps.ByName("username"), ps.ByName("repository"))

		c := pool.Get()
		o, err := c.Do("HGET", k, idfr)
		if err != nil {
			// Log the error and fallback
			lgr.Error(err)
			next.ServeHTTP(w, r)
			return
		}
		if o == nil {
			next.ServeHTTP(w, r)
			return
		}
		lgr.Infoln(k, idfr, "cache hit")
		outputFromCache(w, r, http.StatusOK, o.([]byte))
	}
	return http.HandlerFunc(fn)
}

func setLogger(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		params := context.Get(r, "params").(httprouter.Params)
		registry, username, repository := params.ByName("registry"), params.ByName("username"), params.ByName("repository")

		// Retrieve the user IP
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			handleError(w, r, fmt.Errorf("userip: %q is not IP:port", r.RemoteAddr))
			return
		}

		// Set the logging context
		logger := log.WithFields(log.Fields{
			"registry":   registry,
			"username":   username,
			"repository": repository,
			"ip":         ip,
		})

		context.Set(r, "logger", logger)
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
