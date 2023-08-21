package data

import (
	"database/sql"
	"time"

	"crossfitbox.booking.system/internal/validator"
	"github.com/lib/pq"
)

type WorkoutModel struct {
	DB *sql.DB
}

func (w WorkoutModel) Insert(workout *Workout) error {
	query := `
		INSERT INTO workouts (title, mode, time_cap, equipment, exercises, trainer_tips)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, updated_at, created_at`

	args := []interface{}{
		workout.Title,
		workout.Mode,
		workout.TimeCap,
		pq.Array(workout.Equipment),
		pq.Array(workout.Exercises),
		pq.Array(workout.TrainerTips),
	}

	return w.DB.QueryRow(query, args...).Scan(&workout.ID, &workout.UpdatedAt, &workout.CreatedAt)

}

func (w WorkoutModel) Get(id int64) (*Workout, error) {
	return nil, nil
}

func (w WorkoutModel) Update(workout *Workout) error {
	return nil
}

func (w WorkoutModel) Delete(id int64) error {
	return nil
}

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
