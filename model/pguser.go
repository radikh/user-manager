// Package model provides user-manager specific data structures,
// which are meant to be used across the whole application.
package model

import "database/sql"

//UsersRepo structure that contain pointer to database
type UsersRepo struct {
	db *sql.DB
}

// NeUsersRepo returns UsersRepo with db
func NewUsersRepo(data *sql.DB) *UsersRepo {
	return &UsersRepo{db: data}
}
