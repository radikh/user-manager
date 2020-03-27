package pgclient

import (
	"database/sql"
	"time"

	"github.com/lvl484/user-manager/model"
)

type Postgres struct {
	db *sql.DB
}

// NewPostgres returns Postgres with db
func NewPostgres(data *sql.DB) *Postgres {
	return &Postgres{db: data}
}

const (
	queryInsert     = `INSERT INTO users(id, user_name,password,email,first_name, last_name, phone, created_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`
	queryUpdate     = `UPDATE users SET (password,email,first_name, last_name, phone, updated_at)=($1,$2,$3,$4,$5,$6) WHERE user_name=$7`
	queryDelete     = `DELETE FROM users WHERE user_name=$1`
	queryDisable    = ``
	querySelectInfo = `SELECT user_name,email,first_name, last_name, phone FROM users WHERE user_name=$1`
)

// AddUser adds new user to database
func (pg *Postgres) AddUser(user *model.User) error {
	_, err := pg.db.Query(queryInsert, user.ID, user.Username, user.Password, user.Email, user.FirstName, user.LastName, user.Phone, time.Now())

	return err
}

// UpdateUser update information about user in database
func (pg *Postgres) UpdateUser(user *model.User) error {
	_, err := pg.db.Exec(queryUpdate, user.Password, user.Email, user.FirstName, user.LastName, user.Phone, time.Now(), user.Username)

	return err
}

// DeleteUser delete information about user in database
func (pg *Postgres) DeleteUser(login string) error {
	_, err := pg.db.Exec(queryDelete, login)

	return err
}

// DisableUser deactivate information about user in database
func (pg *Postgres) DisableUser(login string) error {
	_, err := pg.db.Exec(queryDisable, login)

	return err
}

// GetUserInfo get user information from database
func (pg *Postgres) GetUserInfo(login string) (*model.User, error) {
	var usr model.User
	row, err := pg.db.Query(querySelectInfo, login)
	if err != nil {
		return nil, err
	}
	defer row.Close()
	for row.Next() {
		err = row.Scan(usr.Username, usr.Email, usr.FirstName, usr.LastName, usr.Phone)
		if err != nil {
			return nil, err
		}
	}

	return &usr, nil
}
