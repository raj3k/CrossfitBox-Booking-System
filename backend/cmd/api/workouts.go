package main

import (
	"fmt"
	"net/http"
)

func (app *application) createWorkoutHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "create a workout")
}

func (app *application) showWorkoutHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	fmt.Fprintf(w, "show details of workout %d\n", id)
}
