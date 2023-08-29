package data

import (
	"errors"
	"time"

	"crossfitbox.booking.system/internal/types"
	"crossfitbox.booking.system/internal/validator"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
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
}

type password struct {
	plaintext *string
	hash      []byte
}

// The Set() method calculates the bcrypt hash of a plaintext password, and stores both
// the hash and the plaintext versions in the struct
func (p *password) Set(plaintextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 12)
	if err != nil {
		return err
	}

	p.plaintext = &plaintextPassword
	p.hash = hash

	return nil
}

// The Matches() method checks whether the provided plaintext password matches the
// hashed password stored in the struct, returning true if it matches and false
// otherwise.
func (p *password) Matches(plaintextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(plaintextPassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}
	return true, nil
}

func ValidateEmail(v *validator.Validator, email string) {
	v.Check(email != "", "email", "must be provided")
	v.Check(validator.Matches(email, validator.EmailRX), "email", "must be a valid email address")
}

func ValidatePasswordPlaintext(v *validator.Validator, password string) {
	v.Check(password != "", "password", "must be provided")
	v.Check(len(password) >= 8, "password", "must be at least 8 bytes long")
	v.Check(len(password) <= 72, "password", "must not be more than 72 bytes long")
}

func ValidateUser(v *validator.Validator, user *User) {
	v.Check(user.FirstName != "", "first_name", "must be provided")
	v.Check(user.LastName != "", "last_name", "must be provided")

	ValidateEmail(v, user.Email)

	if user.Password.plaintext != nil {
		ValidatePasswordPlaintext(v, *user.Password.plaintext)
	}

	if user.Password.hash == nil {
		panic("missing password hash for user")
	}
}
