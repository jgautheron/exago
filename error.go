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
		outputJSON(w, r, http.StatusBadRequest, err.Error())
	case notFound:
		outputJSON(w, r, http.StatusNotFound, err.Error())
	default:
		outputJSON(w, r, http.StatusInternalServerError, err.Error())
	}
}
