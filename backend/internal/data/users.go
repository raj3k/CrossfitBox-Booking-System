package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"crossfitbox.booking.system/internal/types"
	"crossfitbox.booking.system/internal/validator"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrDuplicateEmail = errors.New("duplicate email")
)

type UserModel struct {
	DB *sql.DB
}

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
	CreatedAt   time.Time   `json:"created_at"`
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

func (um *UserModel) Insert(user *User) error {
	// TODO: return also user Id, will be useful later
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	tx, err := um.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	query_user := `
	INSERT INTO users (email, password, first_name, last_name) 
	VALUES ($1, $2, $3, $4) 
	RETURNING id, created_at`

	args_user := []interface{}{
		user.Email,
		user.Password.hash,
		user.FirstName,
		user.LastName,
	}

	err = tx.QueryRowContext(ctx, query_user, args_user...).Scan(
		&user.ID,
		&user.CreatedAt,
	)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return ErrDuplicateEmail
		default:
			return err
		}
	}

	query_user_profile := `
	INSERT INTO user_profile (user_id)
	VALUES ($1)
	ON CONFLICT (user_id) DO NOTHING RETURNING id, user_id`

	err = tx.QueryRowContext(ctx, query_user_profile, user.ID).Scan(
		&user.Profile.ID,
		&user.Profile.UserID,
	)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (um *UserModel) GetByEmail(email string, active bool) (*User, error) {
	query := `
	SELECT u.*, p.*
	FROM users u
	JOIN user_profile p ON p.user_id = u.id
	WHERE u.is_active = $2 AND u.email = $1`

	var user User
	var userProfile UserProfile

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := um.DB.QueryRowContext(ctx, query, email, active).Scan(
		&user.ID,
		&user.Email,
		&user.Password.hash,
		&user.FirstName,
		&user.LastName,
		&user.IsActive,
		&user.IsStaff,
		&user.IsSuperuser,
		&user.Thumbnail,
		&user.CreatedAt,
		&userProfile.ID,
		&userProfile.UserID,
		&userProfile.PhoneNumber,
		&userProfile.BirthDate,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			if active {
				return nil, ErrRecordNotFound
			} else {
				return nil, errors.New("an inactive user with the provided email address was not found")
			}
		default:
			return nil, err
		}
	}

	user.Profile = userProfile

	return &user, nil
}

func (um *UserModel) Update(user *User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	tx, err := um.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil
	}

	query_user := `
	UPDATE
		users
	SET
		first_name = COALESCE($1, first_name),
		last_name = COALESCE($2, last_name),
		thumbnail = COALESCE($3, thumbnail)
	WHERE
		id = $4 AND is_active = true
	RETURNING
		id,
		email,
		first_name,
		last_name,
		is_active,
		is_staff,
		is_superuser,
		thumbnail,
		created_at`

	args_user := []interface{}{
		user.FirstName,
		user.LastName,
		user.Thumbnail,
		user.ID,
	}

	err = um.DB.QueryRowContext(ctx, query_user, args_user...).Scan(
		&user.ID,
		&user.Email,
		&user.Password.hash,
		&user.FirstName,
		&user.LastName,
		&user.IsActive,
		&user.IsStaff,
		&user.IsSuperuser,
		&user.Thumbnail,
		&user.CreatedAt,
	)

	if err != nil {
		return err
	}

	query_user_profile := `
	UPDATE
		user_profile
	SET
		phone_number = NULLIF($1, ''),
		birth_date = $2::timestamp::date
	WHERE
		user_id = $3
	RETURNING
		id,
		user_id,
		phone_number,
		birth_date`

	args_user_profile := []interface{}{
		user.Profile.PhoneNumber,
		user.Profile.BirthDate,
		user.ID,
	}

	err = tx.QueryRowContext(ctx, query_user_profile, args_user_profile...).Scan(
		&user.Profile.ID,
		&user.Profile.UserID,
		&user.Profile.PhoneNumber,
		&user.Profile.BirthDate,
	)

	if err != nil {
		return err
	}

	return tx.Commit()
}

func (um *UserModel) Activate(userID uuid.UUID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `UPDATE users SET is_active = true WHERE id = $1`

	_, err := um.DB.ExecContext(ctx, query, userID)
	if err != nil {
		return err
	}

	return nil
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
