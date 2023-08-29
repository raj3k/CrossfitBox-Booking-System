package data

import (
	"time"

	"crossfitbox.booking.system/internal/types"
	"github.com/google/uuid"
)

type User struct {
	ID          uuid.UUID   `json:"id"`
	Email       string      `json:"email"`
	Password    password    `json:"-"`
	FirstName   string      `json:"first_name"`
	LastName    string      `json:"last_name"`
	IsActive    bool        `json:"is_active"`
	IsStaff     bool        `json:"is_staff"`
	IsSuperuser bool        `json:"is_superuser"`
	Thumbnail   *string     `json:"thumbnail"`
	DateJoined  time.Time   `json:"date_joined"`
	Profile     UserProfile `json:"profile"`
}

type UserProfile struct {
	ID          *uuid.UUID     `json:"id"`
	UserID      *uuid.UUID     `json:"user_id"`
	PhoneNumber *string        `json:"phone_number"`
	BirthDate   types.NullTime `json:"birth_date"`
	GithubLink  *string        `json:"github_link"`
}

type password struct {
	plaintext *string
	hash      []byte
}
