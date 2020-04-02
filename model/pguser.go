// Command umcli provides admin command line tool to manipulate accounts with admin rights.
package main

import (
	"database/sql"
	//	"errors"
	"time"

	"github.com/pkg/errors"

	"github.com/google/uuid"
	"github.com/lvl484/user-manager/model"
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

// usersRepo structure that contain pointer to database
type usersRepo struct {
	db *sql.DB
}

// NewUsersRepo returns usersRepo with db
func NewUsersRepo(data *sql.DB) *usersRepo {
	return &usersRepo{db: data}
}

// Add adds new user to database
func (ur *usersRepo) Add(user *model.User) error {
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
func (ur *usersRepo) Update(user *model.User) error {
	salted, err := ur.getUserDeactivated(user.Username)
	if err != nil {
		return err
	}
	if salted {
		return ErrUserDisable
	}
	pwd, err := EncodePassword(NewPasswordConfig(), user.Password)
	if err != nil {
		return err
	}
	_, err = ur.db.Exec(queryUpdate, pwd, user.Email, user.FirstName, user.LastName, user.Phone, time.Now(), user.Username)

	return err
}

// Delete delete information about user in database
func (ur *usersRepo) Delete(login string) error {
	_, err := ur.db.Exec(queryDelete, login)

	return err
}

// Disable deactivate information about user in database
func (ur *usersRepo) Disable(login string) error {
	_, err := ur.db.Exec(queryDisable, "true", login)

	return err
}

// Activate deactivate information about user in database
func (ur *usersRepo) Activate(login string) error {
	_, err := ur.db.Exec(queryDisable, "false", login)

	return err
}

// GetInfo get user information from database
func (ur *usersRepo) GetInfo(login string) (*model.User, error) {
	salted, err := ur.getUserDeactivated(login)
	if err != nil {
		return nil, err
	}
	if salted {
		return nil, ErrUserDisable
	}
	var usr model.User
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
func (ur *usersRepo) getUserDeactivated(login string) (bool, error) {
	result := false
	err := ur.db.QueryRow(queryAlive, login).Scan(&result)

	return result, err
}

// CheckLoginExist check information about existing user with such login
func (ur *usersRepo) CheckLoginExist(login string) (bool, error) {
	result := 0
	err := ur.db.QueryRow(queryCheckLogin, login).Scan(&result)

	return result == 1, err
}
