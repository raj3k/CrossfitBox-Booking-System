package data

import (
	"context"
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

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	err := w.DB.QueryRowContext(ctx, query, args...).
		Scan(&workout.ID, &workout.UpdatedAt, &workout.CreatedAt)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "workouts_title_key"`:
			return ErrDuplicateTitle
		default:
			return err
		}
	}
	return nil
}

func (w WorkoutModel) Get(id uuid.UUID) (*Workout, error) {
	query := `
	SELECT id, title, mode, time_cap, equipment, exercises, trainer_tips, created_at, updated_at
	FROM workouts
	WHERE id = $1`

	var workout Workout

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	err := w.DB.QueryRowContext(ctx, query, id).Scan(
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

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	err := w.DB.QueryRowContext(ctx, query, args...).Scan(&workout.ID, &workout.UpdatedAt)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "workouts_title_key"`:
			return ErrDuplicateTitle
		default:
			return err
		}
	}
	return nil
}

func (w WorkoutModel) Delete(id uuid.UUID) error {
	query := `DELETE FROM workouts WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	result, err := w.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}
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
