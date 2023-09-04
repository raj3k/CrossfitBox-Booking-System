package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)

	router.MethodNotAllowed = http.HandlerFunc(app.MethodNotAllowed)

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)

	router.HandlerFunc(http.MethodGet, "/v1/workouts", app.listWorkoutsHandler)
	router.HandlerFunc(http.MethodPost, "/v1/workouts", app.createWorkoutHandler)
	router.HandlerFunc(http.MethodGet, "/v1/workouts/:id", app.showWorkoutHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/workouts/:id", app.updateWorkoutHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/workouts/:id", app.deleteWorkoutHandler)

	router.HandlerFunc(http.MethodPost, "/v1/users/register", app.registerUserHandler)
	router.HandlerFunc(http.MethodPut, "/v1/users/activate/:id/", app.activateUserHandler)

	return app.recoverPanic(router)
}
