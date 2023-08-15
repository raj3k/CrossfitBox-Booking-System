package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() *httprouter.Router {
	// Initialize a new httprouter router instance.
	router := httprouter.New()

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck/", app.healthcheckHandler)
	router.HandlerFunc(http.MethodPost, "/v1/workouts/", app.createWorkoutHandler)
	router.HandlerFunc(http.MethodGet, "/v1/workouts/:id", app.showWorkoutHandler)

	return router
}
