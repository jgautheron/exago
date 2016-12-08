package server

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"

	log "github.com/Sirupsen/logrus"
	. "github.com/hotolab/exago-svc/config"
	"github.com/pressly/chi"
	"github.com/pressly/chi/render"
)

var (
	ErrInvalidVersion     = errors.New("The specified Go version is not supported")
	ErrBranchDoesNotExist = errors.New("The specified branch does not exist")
)

func (s *Server) checkValidData(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			host, owner, name string
			branch, goversion string
		)

		if r.Method == http.MethodPost {
			var data struct {
				Host      string `json:"host"`
				Owner     string `json:"owner"`
				Name      string `json:"name"`
				Branch    string `json:"repository"`
				GoVersion string `json:"goversion"`
			}
			if err := render.Bind(r.Body, &data); err != nil {
				render.Status(r, http.StatusInternalServerError)
				render.JSON(w, r, err.Error())
				return
			}
			host, owner, name = data.Host, data.Owner, data.Name
			branch, goversion = data.Branch, data.GoVersion
		} else {
			host, owner, name = chi.URLParam(r, "host"), chi.URLParam(r, "owner"), chi.URLParam(r, "name")
			branch, goversion = chi.URLParam(r, "branch"), chi.URLParam(r, "goversion")
		}

		repo := fmt.Sprintf("%s/%s/%s", host, owner, name)
		data, err := s.config.RepositoryLoader.IsValid(repo)
		if err != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, err.Error())
			return
		}

		// Validate the branch name
		branchMatch := false
		for _, branchName := range data["branches"].([]string) {
			if branchName == branch {
				branchMatch = true
			}
		}
		if !branchMatch {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, ErrBranchDoesNotExist)
			return
		}

		// Validate the Go version
		versionMatch := false
		for _, allowedVersion := range Config.GoVersions {
			if allowedVersion == goversion {
				versionMatch = true
			}
		}
		if !versionMatch {
			render.Status(r, http.StatusNotImplemented)
			render.JSON(w, r, ErrInvalidVersion)
			return
		}

		// Reuse the proper owner/repository name with caps if any
		HTMLURL := strings.Replace(data["html_url"].(string), "https://", "", 1)
		rp := strings.Split(HTMLURL, "/")

		ctx := context.WithValue(r.Context(), "host", rp[0])
		ctx = context.WithValue(ctx, "owner", rp[1])
		ctx = context.WithValue(ctx, "name", rp[2])
		ctx = context.WithValue(ctx, "repo", HTMLURL)
		ctx = context.WithValue(ctx, "branch", branch)
		ctx = context.WithValue(ctx, "goversion", goversion)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (s *Server) setLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lgr := logger.WithFields(log.Fields{
			"repo":      r.Context().Value("repo").(string),
			"branch":    r.Context().Value("branch").(string),
			"goversion": r.Context().Value("goversion").(string),
			"ip":        getIP(r.RemoteAddr),
		})
		ctx := context.WithValue(r.Context(), "lgr", lgr)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getIP(s string) string {
	if ip, _, err := net.SplitHostPort(s); err == nil {
		return ip
	}
	if ip := net.ParseIP(s); ip != nil {
		return ip.To4().String()
	}
	logger.Errorf("Couldn't parse IP %s", s)
	return ""
}
