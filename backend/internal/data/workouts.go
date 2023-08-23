package data

import (
	"database/sql"
	"errors"
	"time"

	"crossfitbox.booking.system/internal/validator"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type WorkoutModel struct {
	DB *sql.DB
}

func (w WorkoutModel) Insert(workout *Workout) error {
	tx, err := w.DB.Begin()
	if err != nil {
		return err
	}

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

	err = tx.QueryRow(query, args...).
		Scan(&workout.ID, &workout.UpdatedAt, &workout.CreatedAt)
	if err != nil {
		rollbackErr := tx.Rollback()
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "workouts_title_key"`:
			return ErrDuplicateTitle
		case errors.Is(err, rollbackErr):
			return errors.New("rollback error")
		default:
			return err
		}
	}
	return tx.Commit()
}

func (w WorkoutModel) Get(id uuid.UUID) (*Workout, error) {
	query := `
	SELECT id, title, mode, time_cap, equipment, exercises, trainer_tips, created_at, updated_at
	FROM workouts
	WHERE id = $1`

	var workout Workout

	err := w.DB.QueryRow(query, id).Scan(
		&workout.ID,
		&workout.Title,
		&workout.Mode,
		&workout.TimeCap,
		pq.Array(&workout.Equipment),
		pq.Array(&workout.Exercises),
		pq.Array(&workout.TrainerTips),
		&workout.CreatedAt,
		&workout.UpdatedAt,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &workout, nil
}

func (w WorkoutModel) Update(workout *Workout) error {
	query := `
		UPDATE workouts
		SET title = $1, mode = $2, time_cap = $3, equipment = $4, exercises = $5, trainer_tips = $6, updated_at = NOW()
		WHERE id = $7
		RETURNING id, updated_at`
	args := []interface{}{
		workout.Title,
		workout.Mode,
		workout.TimeCap,
		pq.Array(workout.Equipment),
		pq.Array(workout.Exercises),
		pq.Array(workout.TrainerTips),
		workout.ID,
	}
	return w.DB.QueryRow(query, args...).Scan(&workout.ID, &workout.UpdatedAt)
}

func (w WorkoutModel) Delete(id int64) error {
	return nil
}

type Workout struct {
	ID          uuid.UUID `json:"id"`
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
