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
	"github.com/exago/svc/github"
	"github.com/exago/svc/requestlock"
	"github.com/gorilla/context"
	"github.com/julienschmidt/httprouter"
)

const (
	rateLimitCount = 20
)

var (
	ErrInvalidParameter  = errors.New("Invalid parameter passed")
	ErrTooManyCalls      = errors.New("Too many calls in a short period of time")
	ErrRateLimitExceeded = errors.New("Rate limit exceeded")
	ErrInvalidLanguage   = errors.New("The repository does not contain Go code")
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

func checkValidRepository(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ps := context.Get(r, "params").(httprouter.Params)
		repository := ps.ByName("repository")[1:]

		if !strings.HasPrefix(repository, "github.com") {
			writeData(w, r, http.StatusNotImplemented, "Only GitHub is implemented right now")
			return
		}
		re := regexp.MustCompile(`^github\.com/[\w\d\-_]+/[\w\d\-_]+`)
		if !re.MatchString(repository) {
			writeError(w, r, ErrInvalidParameter)
			return
		}

		// Split the project path
		sp := strings.Split(repository, "/")
		context.Set(r, "provider", sp[0])
		context.Set(r, "owner", sp[1])
		context.Set(r, "project", sp[2])

		// Check with the GitHub API if the repository contains Go code
		data, err := github.Get(sp[1], sp[2])
		if err != nil {
			writeError(w, r, err)
			return
		}

		if !strings.Contains(data["language"], "Go") {
			writeError(w, r, ErrInvalidLanguage)
			return
		}

		// Entire path
		context.Set(r, "repository", repository)

		if len(sp) > 3 {
			context.Set(r, "path", strings.Join(sp[3:], "/"))
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
			writeError(w, r, ErrTooManyCalls)
			return
		}

		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func setLogger(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ps := context.Get(r, "params").(httprouter.Params)

		lgr := log.WithFields(log.Fields{
			"repository": ps.ByName("repository")[1:],
			"ip":         getIP(r.RemoteAddr),
		})
		context.Set(r, "logger", lgr)
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
			writeError(w, r, ErrRateLimitExceeded)
			return
		}
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
