package data

import (
	"time"

	"crossfitbox.booking.system/internal/validator"
)

type Workout struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Mode        string    `json:"mode"`
	TimeCap     TimeCap   `json:"time_cap,omitempty"`
	Equipment   []string  `json:"equipment,omitempty"`
	Exercises   []string  `json:"exercises"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	TrainerTips []string  `json:"trainer_tips,omitempty"`
}

func ValidateWorkout(v *validator.Validator, workout *Workout) {
	v.Check(workout.Title != "", "title", "must be provided")
	v.Check(len(workout.Title) <= 500, "title", "must not be more than 500 bytes long")

	v.Check(workout.Mode != "", "mode", "must be provided")

	if workout.TimeCap != 0 {
		v.Check(workout.TimeCap > 0, "time_cap", "must be a positive integer")
	}

	v.Check(workout.Exercises != nil, "exercises", "must be provided")
	v.Check(len(workout.Exercises) >= 1, "exercises", "must contain at least 1 exercise")

	v.Check(validator.Unique(workout.TrainerTips), "trainer_tips", "must not contain duplicate records")
}
