// Package model provides user-manager specific data structures,
// which are meant to be used across the whole application.
package model

import (
	"database/sql"
	//	"errors"
	"time"

	"github.com/pkg/errors"

	"github.com/google/uuid"
)

const (
	queryInsert     = `INSERT INTO users(id, user_name,password,email,first_name, last_name, phone, created_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`
	queryUpdate     = `UPDATE users SET (password,email,first_name, last_name, phone, updated_at)=($1,$2,$3,$4,$5,$6) WHERE user_name=$7`
	queryDelete     = `DELETE FROM users WHERE user_name=$1`
	queryDisable    = `UPDATE users SET salted=$1 WHERE user_name=$2`
	querySelectInfo = `SELECT id,user_name,email,first_name, last_name, phone FROM users WHERE user_name=$1`
	queryAlive      = `SELECT salted FROM users WHERE user_name=$1`
	queryCheckLogin = `SELECT count(id) FROM users WHERE user_name=$1`
)

var errUserDisable = errors.New("User is disabled!")

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
		return err
	}
	errors.Wrap
	ui, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	_, err = ur.db.Query(queryInsert, ui, user.Username, pwd, user.Email, user.FirstName, user.LastName, user.Phone, time.Now())
	return err
}

// Update update information about user in database
func (ur *UsersRepo) Update(user *User) error {
	salted, err := ur.getUserDeactivated(user.Username)
	if err != nil {
		return err
	}
	if salted {
		return errUserDisable
	}
	pwd, err := EncodePassword(NewPasswordConfig(), user.Password)
	if err != nil {
		return err
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
	salted, err := ur.getUserDeactivated(login)
	if err != nil {
		return nil, err
	}
	if salted {
		return nil, errUserDisable
	}
	var usr User
	row, err := ur.db.QueryRow(querySelectInfo, login)
	if err != nil {
		return nil, err
	}
	defer row.Close()
	for row.Next() {
		err = row.Scan(&usr.ID, &usr.Username, &usr.Email, &usr.FirstName, &usr.LastName, &usr.Phone)
		if err != nil {
			return nil, err
		}
	}

	return &usr, nil
}

// getUserDeactivated show if user is deactivated
func (ur *UsersRepo) getUserDeactivated(login string) (bool, error) {
	result := false
	err := ur.db.QueryRow(queryAlive, login).Scan(&result)

	return result, err
}

// CheckLoginExist check information about existing user with such login
func (ur *UsersRepo) CheckLoginExist(login string) (bool, error) {
	result := 0
	err := ur.db.QueryRow(queryCheckLogin, login).Scan(&result)

	return result == 1, err
}
