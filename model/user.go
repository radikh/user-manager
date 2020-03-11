package model

import (
	"time"
)

// Struct with all necessary information about user
type User struct {
	// Unique identifier
	ID string `json:"id"`
	// Unique login
	Username string `json:"user_name"`
	// Strong enough password
	Password string `json:"password"`
	// Valid email
	Email string `json:"email"`
	// Obviously first name
	FirstName string `json:"first_name"`
	// Obviously last name
	LastName string `json:"last_name"`
	// Valid phone
	Phone int `json:"phone"`
	// Time when user was created
	CreatedAt *time.Time `json:"created_at"`
	// Time of last changes made
	UpdatedAt *time.Time `json:"updated_at"`
}
