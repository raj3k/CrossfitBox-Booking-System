package main

import (
	"fmt"
	"net/http"
	"time"

	"crossfitbox.booking.system/internal/data"
	"crossfitbox.booking.system/internal/validator"
)

func (app *application) createWorkoutHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title       string       `json:"title"`
		Mode        string       `json:"mode"`
		TimeCap     data.TimeCap `json:"time_cap"`
		Equipment   []string     `json:"equipment"`
		Exercises   []string     `json:"exercises"`
		TrainerTips []string     `json:"trainer_tips"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	workout := &data.Workout{
		Title:       input.Title,
		Mode:        input.Mode,
		TimeCap:     input.TimeCap,
		Equipment:   input.Equipment,
		Exercises:   input.Exercises,
		TrainerTips: input.TrainerTips,
	}

	v := validator.New()

	if data.ValidateWorkout(v, workout); !v.Valid() {
		app.failedValidationErrors(w, r, v.Errors)
		return
	}

	err = app.models.Workouts.Insert(workout)
	if err != nil {
		// TODO: handle error "pq: duplicate key value violates unique constraint "workouts_title_key""
		// app.serveErrorResponse(w, r, err)
		app.conflictResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/workouts/%d", workout.ID))

	app.writeJSON(w, http.StatusCreated, envelope{"workout": workout}, headers)
	if err != nil {
		app.serveErrorResponse(w, r, err)
	}
}

func (app *application) showWorkoutHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	app.logger.Print(app.models.Workouts.Get(id))

	workout := data.Workout{
		ID:        id,
		Title:     "Tommy V",
		Mode:      "For Time",
		Equipment: []string{"barbell, rope"},
		TimeCap:   14,
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
