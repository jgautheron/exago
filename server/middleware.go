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
	"github.com/gorilla/context"
	. "github.com/hotolab/exago-svc/config"
	"github.com/hotolab/exago-svc/github"
	"github.com/hotolab/exago-svc/requestlock"
	"github.com/julienschmidt/httprouter"
)

const (
	rateLimitCount = 20
	sizeLimit      = 1000000 // bytes
)

var (
	ErrInvalidParameter  = errors.New("Invalid parameter passed")
	ErrTooManyCalls      = errors.New("Too many calls in a short period of time")
	ErrRateLimitExceeded = errors.New("Rate limit exceeded")
	ErrInvalidLanguage   = errors.New("The repository does not contain Go code")
	ErrTooLarge          = errors.New("The repository is too large")
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
		re := regexp.MustCompile(`^github\.com/[\w\d\-\._]+/[\w\d\-\._]+`)
		if !re.MatchString(repository) {
			writeError(w, r, ErrInvalidParameter)
			return
		}

		// Check with the GitHub API if the repository contains Go code
		sp := strings.Split(repository, "/")
		data, err := github.GetInstance().Get(sp[1], sp[2])
		if err != nil {
			writeError(w, r, err)
			return
		}

		languages := data["languages"].(map[string]int)
		size, exists := languages["Go"]
		if !exists {
			writeError(w, r, ErrInvalidLanguage)
			return
		}

		if size > sizeLimit {
			writeError(w, r, ErrTooLarge)
			return
		}

		HTMLURL := strings.Replace(data["html_url"].(string), "https://", "", 1)
		repo := strings.Split(HTMLURL, "/")

		context.Set(r, "provider", repo[0])
		context.Set(r, "owner", repo[1])
		context.Set(r, "project", repo[2])
		context.Set(r, "repository", HTMLURL)

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
