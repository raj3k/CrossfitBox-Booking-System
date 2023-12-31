package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)

	router.MethodNotAllowed = http.HandlerFunc(app.MethodNotAllowed)

	router.HandlerFunc(http.MethodGet, "/api/v1/healthcheck", app.healthcheckHandler)

	// Workout related endpoints
	router.HandlerFunc(http.MethodGet, "/api/v1/workouts", app.listWorkoutsHandler)
	router.HandlerFunc(http.MethodPost, "/api/v1/workouts", app.createWorkoutHandler)
	router.HandlerFunc(http.MethodGet, "/api/v1/workouts/:id", app.showWorkoutHandler)
	router.HandlerFunc(http.MethodPatch, "/api/v1/workouts/:id", app.updateWorkoutHandler)
	router.HandlerFunc(http.MethodDelete, "/api/v1/workouts/:id", app.deleteWorkoutHandler)

	// User related endpoints
	router.HandlerFunc(http.MethodPost, "/api/v1/users/register", app.registerUserHandler)
	router.HandlerFunc(http.MethodPost, "/api/v1/users/login", app.loginUserHandler)
	router.HandlerFunc(http.MethodPut, "/api/v1/users/activate/:id/", app.activateUserHandler)
	router.HandlerFunc(http.MethodGet, "/api/v1/users/current-user", app.currentUserHandler)
	router.HandlerFunc(http.MethodPost, "/api/v1/users/logout", app.logoutUserHandler)

	return app.recoverPanic(app.enableCORS(router))
}
