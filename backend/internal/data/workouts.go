package data

import "time"

type Workout struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Mode        string    `json:"mode,omitempty"`
	Equipment   []string  `json:"equipment"`
	Exercises   []string  `json:"exercises"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	TrainerTips []string  `json:"trainer_tips"`
}
