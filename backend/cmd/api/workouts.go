package main

import (
	"errors"
	"fmt"
	"net/http"

	"crossfitbox.booking.system/internal/data"
	"crossfitbox.booking.system/internal/validator"
)

func (app *application) listWorkoutsHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title     string
		Mode      string
		Equipment []string
		Page      int
		PageSize  int
		Sort      string
	}

	v := validator.New()

	qs := r.URL.Query()

	input.Title = app.readString(qs, "title", "")
	input.Mode = app.readString(qs, "mode", "")
	input.Equipment = app.readCSV(qs, "equipment", []string{})
	input.Page = app.readInt(qs, "page", 1, v)
	input.PageSize = app.readInt(qs, "page_size", 20, v)
	input.Sort = app.readString(qs, "sort", "id")

	if !v.Valid() {
		app.failedValidationErrors(w, r, v.Errors)
		return
	}

	fmt.Fprintf(w, "%+v\n", input)
}

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
		switch {
		case errors.Is(err, data.ErrDuplicateTitle):
			v.AddError("title", "Workout with this title already exists")
			app.failedValidationErrors(w, r, v.Errors)
		default:
			app.serveErrorResponse(w, r, err)
		}
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

	workout, err := app.models.Workouts.Get(*id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serveErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"workout": workout}, nil)
	if err != nil {
		app.serveErrorResponse(w, r, err)
	}
}

func (app *application) updateWorkoutHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	workout, err := app.models.Workouts.Get(*id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serveErrorResponse(w, r, err)
		}
		return
	}

	var input struct {
		Title       *string       `json:"title"`
		Mode        *string       `json:"mode"`
		TimeCap     *data.TimeCap `json:"time_cap"`
		Equipment   []string      `json:"equipment"`
		Exercises   []string      `json:"exercises"`
		TrainerTips []string      `json:"trainer_tips"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.Title != nil {
		workout.Title = *input.Title
	}

	if input.Mode != nil {
		workout.Mode = *input.Mode
	}

	if input.TimeCap != nil {
		workout.TimeCap = *input.TimeCap
	}

	if input.Equipment != nil {
		workout.Equipment = input.Equipment
	}

	if input.Exercises != nil {
		workout.Exercises = input.Exercises
	}

	if input.TrainerTips != nil {
		workout.TrainerTips = input.TrainerTips
	}

	v := validator.New()

	if data.ValidateWorkout(v, workout); !v.Valid() {
		app.failedValidationErrors(w, r, v.Errors)
		return
	}

	err = app.models.Workouts.Update(workout)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateTitle):
			v.AddError("title", "Workout with this title already exists")
			app.failedValidationErrors(w, r, v.Errors)
		default:
			app.serveErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"workout": workout}, nil)
	if err != nil {
		app.serveErrorResponse(w, r, err)
	}
}

func (app *application) deleteWorkoutHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.models.Workouts.Delete(*id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serveErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusNoContent, nil, nil)
	if err != nil {
		app.serveErrorResponse(w, r, err)
	}
}
