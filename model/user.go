package model

import (
	"time"
)

type User struct {
	ID        string    `json:"id"`         // Unique identifier
	Username  string    `json:"user_name"`  // Unique login
	Password  string    `json:"password"`   // Strong enough password
	Email     string    `json:"email"`      // Valid email
	FirstName string    `json:"first_name"` // Obviously first name
	LastName  string    `json:"last_name"`  // Obviously last name
	Phone     int       `json:"phone"`      // Valid phone
	CreatedAt time.Time `json:"created_at"` // Time when user was created
	UpdatedAt time.Time `json:"updated_at"` // Time of last changes made
}
