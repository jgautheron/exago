package server

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/didip/tollbooth"
	"github.com/exago/svc/config"
	"github.com/exago/svc/logger"
	"github.com/exago/svc/redis"
	"github.com/exago/svc/requestlock"
	"github.com/gorilla/context"
	"github.com/julienschmidt/httprouter"
)

const (
	rateLimitCount = 20
)

var (
	errTooManyCalls      = errors.New("Too many calls in a short period of time")
	errRateLimitExceeded = errors.New("Rate limit exceeded")
)

func recoverHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Errorf("panic: %+v", err)
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
			writeData(w, r, http.StatusNotImplemented, fmt.Sprintf("the registry %s is not yet supported", registry))
			return
		}
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func checkValidRepository(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		params := context.Get(r, "params").(httprouter.Params)
		owner, repository := params.ByName("owner"), params.ByName("repository")
		re := regexp.MustCompile(`^[\w\d\-_]+$`)
		if !re.MatchString(owner) || !re.MatchString(repository) {
			writeError(w, r, errInvalidParameter)
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
		k := fmt.Sprintf("%s/%s/%s", ps.ByName("registry"), ps.ByName("owner"), ps.ByName("repository"))

		c := redis.GetConn()
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
		lgr.Debugln(k, idfr, "cache hit")

		// Output straight from the cache
		w.Header().Set("Access-Control-Allow-Origin", config.Get("AllowOrigin"))
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Content-Encoding", "gzip")
		w.Header().Set("X-Cache", "HIT")
		w.WriteHeader(http.StatusOK)

		if _, err := w.Write(o.([]byte)); err != nil {
			writeError(w, r, err)
		}
	}
	return http.HandlerFunc(fn)
}

func requestLock(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ps := context.Get(r, "params").(httprouter.Params)
		rp := fmt.Sprintf("%s/%s/%s", ps.ByName("registry"), ps.ByName("owner"), ps.ByName("repository"))

		log.Info("has", rp, getIP(r.RemoteAddr))
		log.Info(requestlock.Has(rp, getIP(r.RemoteAddr)))
		if requestlock.Has(rp, getIP(r.RemoteAddr)) {
			writeError(w, r, errTooManyCalls)
			return
		}

		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func setLogger(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ps := context.Get(r, "params").(httprouter.Params)
		rp := fmt.Sprintf("%s/%s/%s", ps.ByName("registry"), ps.ByName("owner"), ps.ByName("repository"))
		context.Set(r, "logger", logger.With(rp, getIP(r.RemoteAddr)))
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func rateLimit(next http.Handler) http.Handler {
	limiter := tollbooth.NewLimiter(rateLimitCount, time.Hour*4)
	limiter.Methods = []string{"GET"}

	fn := func(w http.ResponseWriter, r *http.Request) {
		lgr := context.Get(r, "logger").(*log.Entry)

		if r.Header.Get("Origin") == config.Get("AllowOrigin") {
			next.ServeHTTP(w, r)
			return
		}

		httpErr := tollbooth.LimitByRequest(limiter, r)
		if httpErr != nil {
			w.Header().Add("Content-Type", limiter.MessageContentType)
			lgr.WithField("URL", r.URL.String()).Error("Rate limit exceeded")
			writeError(w, r, errRateLimitExceeded)
			return
		}
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}