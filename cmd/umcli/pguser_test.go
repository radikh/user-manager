// Command umcli provides admin command line tool to manipulate accounts with admin rights.
package main

import (
	"database/sql/driver"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	"github.com/lvl484/user-manager/model"
)

func TestNewUsersRepo(t *testing.T) {

	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	userRepo := NewUsersRepo(db)

	assert.Equal(t, db, userRepo.db)
}

func TestAdd(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	timestamp := time.Now()
	mock.ExpectExec(regexp.QuoteMeta(queryInsert)).
		WithArgs("3b60ac82-5e8f-4010-ac99-2344cfa72ce0", "user1", "$argon2id$v=19$m=65536,t=3,p=1$3ep7s6fHN16+6VhygB4KMg$Gb3C1g]", "email1@company.com", "Pedro", "Petrenko", "77777777777", &timestamp).
		WillReturnResult(driver.RowsAffected(1))

	user := model.User{
		ID:        "3b60ac82-5e8f-4010-ac99-2344cfa72ce0",
		Username:  "user1",
		Password:  "$argon2id$v=19$m=65536,t=3,p=1$3ep7s6fHN16+6VhygB4KMg$Gb3C1g]",
		Email:     "email1@company.com",
		FirstName: "Pedro",
		LastName:  "Petrenko",
		Phone:     "77777777777",
		CreatedAt: &timestamp,
	}
	_, err = db.Exec(queryInsert, user.ID, user.Username, user.Password, user.Email, user.FirstName, user.LastName, user.Phone, user.CreatedAt)

	assert.Equal(t, nil, err)
}

func TestUpdate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	timestamp := time.Now()
	mock.ExpectExec(regexp.QuoteMeta(queryUpdate)).
		WithArgs("$argon2id$v=19$m=65536,t=3,p=1$3ep7s6fHN16+6VhygB4KMg$Gb3C1g]", "email1@company.com", "Pedro", "Petrenko", "77777777777", &timestamp, "user1").
		WillReturnResult(driver.RowsAffected(1))

	user := model.User{
		ID:        "3b60ac82-5e8f-4010-ac99-2344cfa72ce0",
		Username:  "user1",
		Password:  "$argon2id$v=19$m=65536,t=3,p=1$3ep7s6fHN16+6VhygB4KMg$Gb3C1g]",
		Email:     "email1@company.com",
		FirstName: "Pedro",
		LastName:  "Petrenko",
		Phone:     "77777777777",
		UpdatedAt: &timestamp,
	}
	_, err = db.Exec(queryUpdate, user.Password, user.Email, user.FirstName, user.LastName, user.Phone, user.UpdatedAt, user.Username)

	assert.Equal(t, nil, err)
}

func TestDelete(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	userRepo := NewUsersRepo(db)

	mock.ExpectExec(regexp.QuoteMeta(queryDelete)).
		WithArgs("user1").
		WillReturnResult(driver.RowsAffected(1))

	err = userRepo.Delete("user1")
	assert.Equal(t, nil, err)
}

func TestDisable(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	userRepo := NewUsersRepo(db)

	mock.ExpectExec(regexp.QuoteMeta(queryDisable)).
		WithArgs("true", "user1").
		WillReturnResult(driver.RowsAffected(1))

	err = userRepo.Disable("user1")
	assert.Equal(t, nil, err)
}

func TestActivate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	userRepo := NewUsersRepo(db)

	mock.ExpectExec(regexp.QuoteMeta(queryDisable)).
		WithArgs("false", "user1").
		WillReturnResult(driver.RowsAffected(1))

	err = userRepo.Activate("user1")
	assert.Equal(t, nil, err)
}

func TestGetInfo(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	userRepo := NewUsersRepo(db)

	rowsSalted := sqlmock.NewRows([]string{"salted"}).
		AddRow("false")

	mock.ExpectQuery(regexp.QuoteMeta(queryAlive)).
		WithArgs("user1").
		WillReturnRows(rowsSalted)

	rowsInfo := sqlmock.NewRows([]string{"id", "user_name", "email", "first_name", "last_name", "phone"}).
		AddRow("3b60ac82-5e8f-4010-ac99-2344cfa72ce0", "user1", "email1@company.com", "Pedro", "Petrenko", "77777777777")

	mock.ExpectQuery(regexp.QuoteMeta(querySelectInfo)).
		WithArgs("user1").
		WillReturnRows(rowsInfo)

	user, err := userRepo.GetInfo("user1")
	if err != nil {
		t.Fatalf("an error '%s' was not expected when executing GetInfo query", err)
	}
	user1 := model.User{
		ID:        "3b60ac82-5e8f-4010-ac99-2344cfa72ce0",
		Username:  "user1",
		Email:     "email1@company.com",
		FirstName: "Pedro",
		LastName:  "Petrenko",
		Phone:     "77777777777",
	}
	assert.Equal(t, &user1, user)

}

func TestGetUserDeactivated(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	userRepo := NewUsersRepo(db)

	rows := sqlmock.NewRows([]string{"salted"}).
		AddRow("false")

	mock.ExpectQuery(regexp.QuoteMeta(queryAlive)).
		WithArgs("user1").
		WillReturnRows(rows)

	isAlive, err := userRepo.getUserDeactivated("user1")
	if err != nil {
		t.Fatalf("an error '%s' was not expected when executing getUserDeactivated query", err)
	}
	assert.Equal(t, false, isAlive)
}

func TestCheckLoginExist(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	userRepo := NewUsersRepo(db)

	rowsExist := sqlmock.NewRows([]string{"count"}).
		AddRow(1)

	mock.ExpectQuery(regexp.QuoteMeta(queryCheckLogin)).
		WithArgs("user1").
		WillReturnRows(rowsExist)

	isExist, err := userRepo.CheckLoginExist("user1")
	if err != nil {
		t.Fatalf("an error '%s' was not expected when executing CheckLoginExist query", err)
	}
	assert.Equal(t, true, isExist)

	rowsNoExist := sqlmock.NewRows([]string{"count"}).
		AddRow(0)

	mock.ExpectQuery(regexp.QuoteMeta(queryCheckLogin)).
		WithArgs("user1").
		WillReturnRows(rowsNoExist)

	isAbsent, err := userRepo.CheckLoginExist("user1")
	if err != nil {
		t.Fatalf("an error '%s' was not expected when executing CheckLoginExist query", err)
	}
	assert.Equal(t, false, isAbsent)
}
