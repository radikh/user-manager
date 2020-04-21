// Package model provides user-manager specific data structures,
// which are meant to be used across the whole application.
package model

import (
	"database/sql/driver"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
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
	conf := NewPasswordConfig()
	assert.NotNil(t, conf)

	pwd, err := EncodePassword(conf, "password")
	if err != nil {
		assert.Error(t, errors.Wrap(err, msgErrorHashingPassword))
	}
	assert.NotNil(t, pwd)
	assert.NoError(t, err)

	ui, err := uuid.NewRandom()
	if err != nil {
		assert.Error(t, errors.Wrap(err, msgErrorGeneratingUUID))
	}
	assert.NotNil(t, ui)
	assert.NoError(t, err)

	timestamp := time.Now()
	mock.ExpectExec(regexp.QuoteMeta(queryInsert)).
		WithArgs("3b60ac82-5e8f-4010-ac99-2344cfa72ce0", "user1", "$argon2id$v=19$m=65536,t=3,p=1$3ep7s6fHN16+6VhygB4KMg$Gb3C1g]", "email1@company.com", "Pedro", "Petrenko", "77777777777", &timestamp).
		WillReturnResult(driver.RowsAffected(1))

	user := User{
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

	assert.NoError(t, err)
}

func TestUpdate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	conf := NewPasswordConfig()
	assert.NotNil(t, conf)

	pwd, err := EncodePassword(conf, "password")
	if err != nil {
		assert.Error(t, errors.Wrap(err, msgErrorHashingPassword))
	}
	assert.NotNil(t, pwd)
	assert.NoError(t, err)

	timestamp := time.Now()
	user := User{
		ID:        "3b60ac82-5e8f-4010-ac99-2344cfa72ce0",
		Username:  "user1",
		Password:  "$argon2id$v=19$m=65536,t=3,p=1$OLBzWepNZEtV3LXyp7SuHQ$tN8Q03tH+lEjUDuxJ1vX+w]",
		Email:     "email1@company.com",
		FirstName: "Pedro",
		LastName:  "Petrenko",
		Phone:     "77777777777",
		UpdatedAt: &timestamp,
	}

	mock.ExpectExec(regexp.QuoteMeta(queryUpdate)).
		WithArgs("$argon2id$v=19$m=65536,t=3,p=1$OLBzWepNZEtV3LXyp7SuHQ$tN8Q03tH+lEjUDuxJ1vX+w]", "email1@company.com", "Pedro", "Petrenko", "77777777777", &timestamp, "user1").
		WillReturnResult(driver.RowsAffected(1))

	_, err = db.Exec(queryUpdate, user.Password, user.Email, user.FirstName, user.LastName, user.Phone, user.UpdatedAt, user.Username)

	assert.NoError(t, err)
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

	mock.ExpectExec(regexp.QuoteMeta(queryDeleteActivationCode)).
		WithArgs("user1").
		WillReturnResult(driver.RowsAffected(1))

	err = userRepo.Delete("user1")
	assert.NoError(t, err)
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
	assert.NoError(t, err)
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
	assert.NoError(t, err)
}

func TestGetInfo(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	user1 := User{
		ID:        "3b60ac82-5e8f-4010-ac99-2344cfa72ce0",
		Username:  "user1",
		Password:  "$argon2id$v=19$m=65536,t=3,p=1$BCDndJ1kUOAAW/mwP7ViOQ$Ig4hpteBW1YM7Lrh3EHkHQ",
		Email:     "email1@company.com",
		FirstName: "Pedro",
		LastName:  "Petrenko",
		Phone:     "77777777777",
	}

	userRepo := NewUsersRepo(db)

	rowsInfo := sqlmock.NewRows([]string{"id", "user_name", "password", "email", "first_name", "last_name", "phone", "salted"}).
		AddRow("3b60ac82-5e8f-4010-ac99-2344cfa72ce0", "user1", "$argon2id$v=19$m=65536,t=3,p=1$BCDndJ1kUOAAW/mwP7ViOQ$Ig4hpteBW1YM7Lrh3EHkHQ",
			"email1@company.com", "Pedro", "Petrenko", "77777777777", "false")

	mock.ExpectQuery(regexp.QuoteMeta(querySelectInfo)).
		WithArgs("user1").
		WillReturnRows(rowsInfo)

	user, err := userRepo.GetInfo("user1")
	assert.NoError(t, err)
	assert.Equal(t, &user1, user)

	rowsDisabled := sqlmock.NewRows([]string{"id", "user_name", "password", "email", "first_name", "last_name", "phone", "salted"}).
		AddRow("3b60ac82-5e8f-4010-ac99-2344cfa72ce0", "user1", "$argon2id$v=19$m=65536,t=3,p=1$BCDndJ1kUOAAW/mwP7ViOQ$Ig4hpteBW1YM7Lrh3EHkHQ",
			"email1@company.com", "Pedro", "Petrenko", "77777777777", "true")

	mock.ExpectQuery(regexp.QuoteMeta(querySelectInfo)).
		WithArgs("user1").
		WillReturnRows(rowsDisabled)

	_, err = userRepo.GetInfo("user1")
	assert.Error(t, err)

}

func TestGetUserInfoIncludingSalted(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	user1 := User{
		ID:        "3b60ac82-5e8f-4010-ac99-2344cfa72ce0",
		Username:  "user1",
		Password:  "$argon2id$v=19$m=65536,t=3,p=1$BCDndJ1kUOAAW/mwP7ViOQ$Ig4hpteBW1YM7Lrh3EHkHQ",
		Email:     "email1@company.com",
		FirstName: "Pedro",
		LastName:  "Petrenko",
		Phone:     "77777777777",
	}

	userRepo := NewUsersRepo(db)

	rowsInfo := sqlmock.NewRows([]string{"id", "user_name", "password", "email", "first_name", "last_name", "phone", "salted"}).
		AddRow("3b60ac82-5e8f-4010-ac99-2344cfa72ce0", "user1", "$argon2id$v=19$m=65536,t=3,p=1$BCDndJ1kUOAAW/mwP7ViOQ$Ig4hpteBW1YM7Lrh3EHkHQ",
			"email1@company.com", "Pedro", "Petrenko", "77777777777", "true")

	mock.ExpectQuery(regexp.QuoteMeta(querySelectInfo)).
		WithArgs("user1").
		WillReturnRows(rowsInfo)

	user, err := userRepo.GetUserInfoIncludingSalted("user1")
	assert.NoError(t, err)
	assert.Equal(t, &user1, user)

	rowsActivated := sqlmock.NewRows([]string{"id", "user_name", "password", "email", "first_name", "last_name", "phone", "salted"}).
		AddRow("3b60ac82-5e8f-4010-ac99-2344cfa72ce0", "user1", "$argon2id$v=19$m=65536,t=3,p=1$BCDndJ1kUOAAW/mwP7ViOQ$Ig4hpteBW1YM7Lrh3EHkHQ",
			"email1@company.com", "Pedro", "Petrenko", "77777777777", "true")

	mock.ExpectQuery(regexp.QuoteMeta(querySelectInfo)).
		WithArgs("user1").
		WillReturnRows(rowsActivated)

	_, err = userRepo.GetUserInfoIncludingSalted("user1")
	assert.NoError(t, err)

	rowsVerified := sqlmock.NewRows([]string{"id", "user_name", "password", "email", "first_name", "last_name", "phone", "salted"}).
		AddRow("3b60ac82-5e8f-4010-ac99-2344cfa72ce0", "user1", "$argon2id$v=19$m=65536,t=3,p=1$BCDndJ1kUOAAW/mwP7ViOQ$Ig4hpteBW1YM7Lrh3EHkHQ",
			"email1@company.com", "Pedro", "Petrenko", "77777777777", "false")

	mock.ExpectQuery(regexp.QuoteMeta(querySelectInfo)).
		WithArgs("user1").
		WillReturnRows(rowsVerified)

	_, err = userRepo.GetUserInfoIncludingSalted("user1")
	assert.Error(t, err)

	assert.EqualError(t, err, "account is already verified")
}
