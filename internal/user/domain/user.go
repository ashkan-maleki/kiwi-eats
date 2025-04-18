package domain

import (
	"errors"
	"time"
)

type User struct {
	ID             string    `json:"id"`
	Email          string    `json:"email"` // Should be validated
	Phone          string    `json:"phone"` // Add for food delivery
	Name           string    `json:"name"`
	HashedPassword []byte    `json:"-"`             // Remove from JSON output
	AvatarURL      string    `json:"avatar_url"`    // Profile picture
	IsVerified     bool      `json:"is_verified"`   // Email/phone verification
	LastLoginAt    time.Time `json:"last_login_at"` // Security tracking
	CreatedAt      time.Time `json:"created_at"`    // Consistent snake_case
	UpdatedAt      time.Time `json:"updated_at"`    // Track modifications
}

var (
	ErrInvalidEmailFormat = errors.New("invalid email format")
)

type EmailValidationFunc func(string) error
type GenerateHashFromPasswordFunc func(password []byte) ([]byte, error)

func NewUser(ID string, email string, name string, password string,
	emailValidationFunc EmailValidationFunc,
	passwordFunc GenerateHashFromPasswordFunc,
	cost int) (*User, error) {
	if err := emailValidationFunc(email); err != nil {
		return nil, errors.Join(ErrInvalidEmailFormat, err)
	}
	hashedPassword, err := passwordFunc([]byte(password))
	if err != nil {
		return nil, err
	}
	return &User{ID: ID, Email: email, Name: name, HashedPassword: hashedPassword}, nil
}
