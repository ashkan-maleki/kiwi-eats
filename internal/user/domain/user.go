package domain

import "time"

type User struct {
	ID             string    `json:"id"`
	Email          string    `json:"email"` // Should be validated
	Phone          string    `json:"phone"` // Add for food delivery
	Name           string    `json:"name"`
	HashedPassword string    `json:"-"`             // Remove from JSON output
	AvatarURL      string    `json:"avatar_url"`    // Profile picture
	IsVerified     bool      `json:"is_verified"`   // Email/phone verification
	LastLoginAt    time.Time `json:"last_login_at"` // Security tracking
	CreatedAt      time.Time `json:"created_at"`    // Consistent snake_case
	UpdatedAt      time.Time `json:"updated_at"`    // Track modifications
}
