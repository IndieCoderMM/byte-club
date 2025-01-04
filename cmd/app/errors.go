package main

import (
	"log"
	"net/http"
)

func (app *application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("Internal server error: %s %s \nError: %s", r.Method, r.URL.Path, err)

	writeJSONError(w, http.StatusInternalServerError, "Internal server error")
}

func (app *application) notFound(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("Not found: %s %s \nError: %s", r.Method, r.URL.Path, err)

	writeJSONError(w, http.StatusNotFound, "Resource not found")
}

func (app *application) badRequest(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("Bad request: %s %s \nError: %s", r.Method, r.URL.Path, err)

	writeJSONError(w, http.StatusBadRequest, err.Error())
}
