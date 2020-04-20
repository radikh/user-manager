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
	queryUpdate = `UPDATE users SET (email,first_name, last_name, phone, 
		updated_at)=($1,$2,$3,$4,$5) WHERE user_name=$6`
	queryDelete     = `DELETE FROM users WHERE user_name=$1`
	queryDisable    = `UPDATE users SET salted=$1 WHERE user_name=$2`
	querySelectInfo = `SELECT id,user_name,password,email,first_name, last_name, 
		phone, salted FROM users WHERE user_name=$1`
	queryUpdatePassword     = `UPDATE users SET (password,updated_at,salted)=($1,$2,false) WHERE user_name=$3`
	queryGetEmail           = `SELECT email FROM users WHERE user_name=$1`
	querySetCode            = `INSERT INTO codes(user_name,code, expired_at, active) VALUES ($1,$2,$3,$4)`
	queryDisableCode        = `UPDATE codes SET active = false WHERE active=true AND user_name=$1`
	queryCheckCode          = `SELECT code, expired_at FROM codes WHERE user_name=$1 AND active=true`
	msgUserDisable          = "User is disabled"
	msgErrorHashingPassword = "Error hashing password"
	msgErrorGeneratingUUID  = "Error generating new UUID for user"
	msgUserDidNotExist      = "There is no such user in database"
	msgCodeDidNotExist      = "There is no any generating code for that user in database"
	msgCodeExpired          = "Activation code is expired"
)

//usersRepo structure that contain pointer to database
type usersRepo struct {
	db *sql.DB
}

// NeusersRepo returns usersRepo with db
func SetUsersRepo(data *sql.DB) *usersRepo {
	return &usersRepo{db: data}
}

// Add adds new user to database
func (ur *usersRepo) Add(user *User) error {
	pwd, err := EncodePassword(NewPasswordConfig(), user.Password)
	if err != nil {
		return errors.Wrap(err, msgErrorHashingPassword)
	}

	ui, err := uuid.NewRandom()
	if err != nil {
		return errors.Wrap(err, msgErrorGeneratingUUID)
	}
	_, err = ur.db.Exec(queryInsert, ui, user.Username, pwd, user.Email, user.FirstName, user.LastName, user.Phone, time.Now())
	return err
}

// Update update information about user in database
func (ur *usersRepo) Update(user *User) error {
	_, err := ur.db.Exec(queryUpdate, user.Email, user.FirstName, user.LastName, user.Phone, time.Now(), user.Username)

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
func (ur *usersRepo) GetInfo(login string) (*User, error) {
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

//UpdatePassword update password for current user account
func (ur *usersRepo) UpdatePassword(login string, password string) error {
	pwd, err := EncodePassword(NewPasswordConfig(), password)
	if err != nil {
		return errors.Wrap(err, msgErrorHashingPassword)
	}
	_, err = ur.db.Exec(queryUpdatePassword, pwd, time.Now(), login)

	return err
}

//GetEmail returnns email for current user
func (ur *usersRepo) GetEmail(login string) (string, error) {
	var email string
	err := ur.db.QueryRow(queryGetEmail, login).Scan(&email)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", errors.Wrap(err, msgUserDidNotExist)
		}
		return "", err
	}
	return email, nil
}

//SetActivationCode write activation code for changing password to database
func (ur *usersRepo) SetActivationCode(login string, code string) error {
	trans, err := ur.db.Begin()
	if err != nil {
		return errors.Wrap(err, trans.Rollback().Error())
	}
	err = ur.Disable(login)
	if err != nil {
		return errors.Wrap(err, trans.Rollback().Error())
	}
	_, err = ur.db.Exec(queryDisableCode, login)
	if err != nil {
		return errors.Wrap(err, trans.Rollback().Error())
	}
	nowTime := time.Now()
	_, err = ur.db.Exec(querySetCode, login, code, nowTime.Add(time.Hour*24), true)
	if err != nil {
		return errors.Wrap(err, trans.Rollback().Error())
	}
	err = trans.Commit()
	return err
}

//CheckActivationCode read from database  activation code for changing password
// and compare it with provided code and return true if match, else return false
func (ur *usersRepo) CheckActivationCode(login string, code string) (bool, error) {
	var dbcode string
	var expTime time.Time
	err := ur.db.QueryRow(queryCheckCode, login).Scan(&dbcode, &expTime)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, errors.Wrap(err, msgCodeDidNotExist)
		}
		return false, err
	}
	if expTime.Before(time.Now()) {
		return false, errors.Wrap(err, msgCodeExpired)
	}
	return dbcode == code, nil
}
