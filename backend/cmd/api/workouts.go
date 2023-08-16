package main

import (
	"fmt"
	"net/http"
	"time"

	"crossfitbox.booking.system/internal/data"
)

func (app *application) createWorkoutHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "create a workout")
}

func (app *application) showWorkoutHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	workout := data.Workout{
		ID:        id,
		CreatedAt: time.Now(),
		Title:     "Fran",
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"workout": workout}, nil)
	if err != nil {
		app.serveErrorResponse(w, r, err)
	}
}
