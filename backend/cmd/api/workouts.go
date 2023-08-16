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
		Title:     "Tommy V",
		Mode:      "For Time",
		Equipment: []string{"barbell, rope"},
		Exercises: []string{
			"21 thrusters",
			"12 rope climbs, 15 ft",
			"15 thrusters",
			"9 rope climbs, 15 ft",
			"9 thrusters",
			"6 rope climbs, 15 ft",
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		TrainerTips: []string{
			"Split the 21 thrusters as needed",
			"Try to do the 9 and 6 thrusters unbroken",
			"RX Weights: 115lb/75lb",
		},
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"workout": workout}, nil)
	if err != nil {
		app.serveErrorResponse(w, r, err)
	}
}
