package model

import (
	"database/sql/driver"
	"regexp"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestDeleteVerificationCode(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	userRepo := NewUsersRepo(db)

	mock.ExpectExec(regexp.QuoteMeta(queryDeleteActivationCode)).
		WithArgs("user1").
		WillReturnResult(driver.RowsAffected(1))

	err = userRepo.DeleteVerificationCode("user1")
	assert.NoError(t, err)
}

func TestGetVerificationCodeTimeSuccess(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	userRepo := NewUsersRepo(db)

	var (
		login       = "bodja"
		code        = "123456789qwerty"
		createdTime = time.Now()
	)

	mock.ExpectQuery(regexp.QuoteMeta(querySelectVerificationCodeTime)).
		WithArgs(login).
		WillReturnRows(sqlmock.NewRows([]string{"code", "createdTime"}).
			AddRow(code, createdTime))

	getTime, getCode, err := userRepo.GetVerificationCodeTime(login)

	require.NoError(t, err)
	assert.NotEmpty(t, getTime)
	assert.NotEmpty(t, getCode)
}

func TestGetVerificationCodeTimeFail(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	userRepo := NewUsersRepo(db)

	var (
		login       = "bodja"
		code        = "123456789qwerty"
		createdTime = time.Now().String()
	)

	mock.ExpectQuery(regexp.QuoteMeta(querySelectVerificationCodeTime)).
		WithArgs(login).
		WillReturnRows(sqlmock.NewRows([]string{"code", "createdTime"}).
			AddRow(code, createdTime))

	_, _, err = userRepo.GetVerificationCodeTime(login)

	require.Error(t, err)
}
