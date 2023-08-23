package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrDuplicateTitle = errors.New("duplicate title")
)

type Models struct {
	Workouts WorkoutModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Workouts: WorkoutModel{DB: db},
	}
}
