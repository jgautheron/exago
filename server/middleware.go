package server

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/didip/tollbooth"
	"github.com/gorilla/context"
	. "github.com/hotolab/exago-svc/config"
	"github.com/hotolab/exago-svc/repository"
	"github.com/hotolab/exago-svc/requestlock"
	"github.com/julienschmidt/httprouter"
)

const (
	rateLimitCount = 20
)

var (
	ErrTooManyCalls      = errors.New("Too many calls in a short period of time")
	ErrRateLimitExceeded = errors.New("Rate limit exceeded")
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
		repo := ps.ByName("repository")[1:]

		data, err := repository.IsValid(repo)
		if err != nil {
			writeError(w, r, err)
			return
		}

		HTMLURL := strings.Replace(data["html_url"].(string), "https://", "", 1)
		rp := strings.Split(HTMLURL, "/")

		context.Set(r, "provider", repo[0])
		context.Set(r, "owner", repo[1])
		context.Set(r, "project", repo[2])
		context.Set(r, "repository", HTMLURL)

		if len(rp) > 3 {
			context.Set(r, "path", strings.Join(rp[3:], "/"))
		}

		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func requestLock(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ps := context.Get(r, "params").(httprouter.Params)
		rp := fmt.Sprintf("%s/%s/%s", ps.ByName("registry"), ps.ByName("owner"), ps.ByName("repository"))

		if requestlock.Contains(rp, getIP(r.RemoteAddr)) {
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

		lgr := logger.WithFields(log.Fields{
			"repository": ps.ByName("repository")[1:],
			"ip":         getIP(r.RemoteAddr),
		})
		context.Set(r, "lgr", lgr)
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func rateLimit(next http.Handler) http.Handler {
	limiter := tollbooth.NewLimiter(rateLimitCount, time.Hour*4)
	limiter.Methods = []string{"GET"}

	fn := func(w http.ResponseWriter, r *http.Request) {
		lgr := context.Get(r, "lgr").(*log.Entry)

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
