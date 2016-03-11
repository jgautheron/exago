package server

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/didip/tollbooth"
	. "github.com/exago/svc/config"
	"github.com/exago/svc/leveldb"
	"github.com/exago/svc/logger"
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
		if !strings.HasPrefix(params.ByName("repository")[1:], "github.com") {
			writeData(w, r, http.StatusNotImplemented, "Only GitHub is implemented right now")
			return
		}
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func initDB(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		_ = leveldb.Instance()
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func checkValidRepository(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ps := context.Get(r, "params").(httprouter.Params)
		re := regexp.MustCompile(`^github\.com/[\w\d\-_]+/[\w\d\-_]+$`)
		if !re.MatchString(ps.ByName("repository")[1:]) {
			writeError(w, r, errInvalidParameter)
			return
		}
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func requestLock(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ps := context.Get(r, "params").(httprouter.Params)
		rp := fmt.Sprintf("%s/%s/%s", ps.ByName("registry"), ps.ByName("owner"), ps.ByName("repository"))

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
		context.Set(r, "logger", logger.With(ps.ByName("repository")[1:], getIP(r.RemoteAddr)))
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func rateLimit(next http.Handler) http.Handler {
	limiter := tollbooth.NewLimiter(rateLimitCount, time.Hour*4)
	limiter.Methods = []string{"GET"}

	fn := func(w http.ResponseWriter, r *http.Request) {
		lgr := context.Get(r, "logger").(*log.Entry)

		if r.Header.Get("Origin") == Config.AllowOrigin {
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
