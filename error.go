package main

import (
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/context"
)

// badRequest is handled by setting the status code in the reply to StatusBadRequest.
type badRequest struct{ error }

// notFound is handled by setting the status code in the reply to StatusNotFound.
type notFound struct{ error }

func handleError(w http.ResponseWriter, r *http.Request, err error) {
	logger := context.Get(r, "logger").(*log.Entry)
	logger.Error(err)

	switch err.(type) {
	case badRequest:
		handleOutput(w, http.StatusBadRequest, err.Error())
	case notFound:
		handleOutput(w, http.StatusNotFound, err.Error())
	default:
		handleOutput(w, http.StatusInternalServerError, err.Error())
	}
}
