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
	queryCheckAdmin         = `SELECT password FROM admins WHERE admin=$1`
	msgUserDisable          = "User is disabled"
	msgErrorHashingPassword = "Error hashing password"
	msgErrorGeneratingUUID  = "Error generating new UUID for user"
	msgUserDidNotExist      = "There is no such user in database"
	msgCodeDidNotExist      = "There is no any generating code for that user in database"
	msgCodeExpired          = "Activation code is expired"
	msgErrorCheckLogin      = "Error: no such user"
	msgErrorCheckPwd        = "Error: password mismatch"
)

//UsersRepo structure that contain pointer to database
type UsersRepo struct {
	db *sql.DB
}

// NeUsersRepo returns UsersRepo with db
func SetUsersRepo(data *sql.DB) *UsersRepo {
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
	_, err = ur.db.Exec(queryInsert, ui, user.Username, pwd, user.Email, user.FirstName, user.LastName, user.Phone, time.Now())
	return err
}

// Update update information about user in database
func (ur *UsersRepo) Update(user *User) error {
	_, err := ur.db.Exec(queryUpdate, user.Email, user.FirstName, user.LastName, user.Phone, time.Now(), user.Username)

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

//UpdatePassword update password for current user account
func (ur *UsersRepo) UpdatePassword(login string, password string) error {
	pwd, err := EncodePassword(NewPasswordConfig(), password)
	if err != nil {
		return errors.Wrap(err, msgErrorHashingPassword)
	}
	_, err = ur.db.Exec(queryUpdatePassword, pwd, time.Now(), login)

	return err
}

//GetEmail returnns email for current user
func (ur *UsersRepo) GetEmail(login string) (string, error) {
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
func (ur *UsersRepo) SetActivationCode(login string, code string) error {
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
func (ur *UsersRepo) CheckActivationCode(login string, code string) (bool, error) {
	var dbcode string
	var expTime time.Time
	err := ur.db.QueryRow(queryCheckCode, login).Scan(&dbcode, &expTime)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, errors.Wrap(err, msgErrorCheckLogin)
		}
		return false, err
	}
	if expTime.Before(time.Now()) {
		return false, errors.Wrap(err, msgCodeExpired)
	}
	return dbcode == code, nil
}

//CheckActivationCode read from database  activation code for changing password
// and compare it with provided code and return true if match, else return false
func (ur *UsersRepo) CheckAdminRole(login string, pwd string) (bool, error) {
	var password string
	err := ur.db.QueryRow(queryCheckAdmin, login).Scan(&password)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, errors.Wrap(err, msgErrorCheckLogin)
		}
		return false, err
	}
	status, err := ComparePassword(pwd, password)
	if err != nil {
		return false, err
	}
	if !status {
		return false, errors.Wrap(err, msgErrorCheckPwd)
	}
	return status, nil
}
