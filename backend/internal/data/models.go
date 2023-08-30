package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrDuplicateName  = errors.New("duplicate name")
)

type Models struct {
	Workouts WorkoutModel
	User     UserModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Workouts: WorkoutModel{DB: db},
		User:     UserModel{DB: db},
	}
}
