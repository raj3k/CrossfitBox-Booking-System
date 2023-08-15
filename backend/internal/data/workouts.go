package data

import "time"

// TODO: design Workout
type Workout struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Title     string    `json:"title"`
	Mode      string    `json:"mode,omitempty"`
}
