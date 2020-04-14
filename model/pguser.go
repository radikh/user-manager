// Package model provides user-manager specific data structures,
// which are meant to be used across the whole application.
package model

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

const (
	queryInsert = `INSERT INTO users(id, user_name,password,email,first_name, 
		last_name, phone, created_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`
	queryUpdate = `UPDATE users SET (password,email,first_name, last_name, phone, 
		updated_at)=($1,$2,$3,$4,$5,$6) WHERE user_name=$7`
	queryDelete             = `DELETE FROM users WHERE user_name=$1`
	queryDisable            = `UPDATE users SET salted=$1 WHERE user_name=$2`
	querySelectInfo         = `SELECT id,user_name,password,email,first_name, last_name, phone, salted FROM users WHERE user_name=$1`
	msgUserDisable          = "User is disabled"
	msgErrorHashingPassword = "Error hashing password"
	msgErrorGeneratingUUID  = "Error generating new UUID for user"
	msgUserDidNotExist      = "There is no such user in database"
)

//UsersRepo structure that contain pointer to database
type UsersRepo struct {
	db *sql.DB
}

// NeUsersRepo returns UsersRepo with db
func NewUsersRepo(data *sql.DB) *UsersRepo {
	return &UsersRepo{db: data}
}

// Add adds new user to database
func (ur *UsersRepo) Add(user *User) error {
	pwd, err := EncodePassword(NewPasswordConfig(), user.Password)
	if err != nil {
		return errors.Wrap(err, msgErrorHashingPassword)
	}

	ui, err := uuid.NewRandom()
	if err != nil {
		return errors.Wrap(err, msgErrorGeneratingUUID)
	}
	user.ID = ui.String()
	_, err = ur.db.Exec(queryInsert, ui, user.Username, pwd, user.Email, user.FirstName, user.LastName, user.Phone, time.Now())
	return err
}

// Update update information about user in database
func (ur *UsersRepo) Update(user *User) error {
	pwd, err := EncodePassword(NewPasswordConfig(), user.Password)
	if err != nil {
		return errors.Wrap(err, msgErrorHashingPassword)
	}
	_, err = ur.db.Exec(queryUpdate, pwd, user.Email, user.FirstName, user.LastName, user.Phone, time.Now(), user.Username)

	return err
}

// Delete delete information about user in database
func (ur *UsersRepo) Delete(login string) error {
	_, err := ur.db.Exec(queryDelete, login)

	return err
}

// Disable deactivate information about user in database
func (ur *UsersRepo) Disable(login string) error {
	_, err := ur.db.Exec(queryDisable, "true", login)

	return err
}

// Activate deactivate information about user in database
func (ur *UsersRepo) Activate(login string) error {
	_, err := ur.db.Exec(queryDisable, "false", login)

	return err
}

// GetInfo get user information from database
func (ur *UsersRepo) GetInfo(login string) (*User, error) {
	var usr User
	var salted bool
	err := ur.db.QueryRow(querySelectInfo, login).Scan(&usr.ID, &usr.Username, &usr.Password,
		&usr.Email, &usr.FirstName, &usr.LastName, &usr.Phone, &salted)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.Wrap(err, msgUserDidNotExist)
		}
		return nil, err
	}
	if salted {
		return nil, errors.New(msgUserDisable)
	}

	return &usr, nil
}
