// Package model provides user-manager specific data structures,
// which are meant to be used across the whole application.
package model

import (
	"database/sql"
	"database/sql/driver"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestSetUsersRepo(t *testing.T) {

	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	userRepo := SetUsersRepo(db)

	assert.Equal(t, db, userRepo.db)
}

func TestAdd(t *testing.T) {
	mock, _, db := mockUserRepo(t)
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
		WithArgs("3b60ac82-5e8f-4010-ac99-2344cfa72ce0", "user1",
			"$argon2id$v=19$m=65536,t=3,p=1$BCDndJ1kUOAAW/mwP7ViOQ$Ig4hpteBW1YM7Lrh3EHkHQ",
			"email1@company.com", "Pedro", "Petrenko", "77777777777", &timestamp).
		WillReturnResult(driver.RowsAffected(1))

	user := mockUser()
	user.CreatedAt = &timestamp
	_, err = db.Exec(queryInsert, user.ID, user.Username, user.Password, user.Email, user.FirstName, user.LastName, user.Phone, user.CreatedAt)

	assert.NoError(t, err)
}

func TestUpdate(t *testing.T) {
	mock, _, db := mockUserRepo(t)
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
	user := mockUser()
	user.UpdatedAt = &timestamp

	mock.ExpectExec(regexp.QuoteMeta(queryUpdate)).
		WithArgs("email1@company.com", "Pedro", "Petrenko", "77777777777", &timestamp, "user1").
		WillReturnResult(driver.RowsAffected(1))

	_, err = db.Exec(queryUpdate, user.Email, user.FirstName, user.LastName, user.Phone, user.UpdatedAt, user.Username)
	assert.NoError(t, err)
}

func TestDelete(t *testing.T) {
	mock, userRepo, db := mockUserRepo(t)
	defer db.Close()
	mock.ExpectExec(regexp.QuoteMeta(queryDelete)).
		WithArgs("user1").
		WillReturnResult(driver.RowsAffected(1))

	err := userRepo.Delete("user1")
	assert.NoError(t, err)
}

func TestDisable(t *testing.T) {
	mock, userRepo, db := mockUserRepo(t)
	defer db.Close()
	mock.ExpectExec(regexp.QuoteMeta(queryDisable)).
		WithArgs("true", "user1").
		WillReturnResult(driver.RowsAffected(1))

	err := userRepo.Disable("user1")
	assert.NoError(t, err)
}

func TestActivate(t *testing.T) {
	mock, userRepo, db := mockUserRepo(t)
	defer db.Close()
	mock.ExpectExec(regexp.QuoteMeta(queryDisable)).
		WithArgs("false", "user1").
		WillReturnResult(driver.RowsAffected(1))

	err := userRepo.Activate("user1")
	assert.NoError(t, err)
}

func TestGetInfo(t *testing.T) {
	mock, userRepo, db := mockUserRepo(t)
	defer db.Close()
	user1 := mockUser()
	str := []string{"id", "user_name", "password", "email", "first_name", "last_name", "phone", "salted"}
	rowsInfo := sqlmock.NewRows(str).
		AddRow("3b60ac82-5e8f-4010-ac99-2344cfa72ce0", "user1",
			"$argon2id$v=19$m=65536,t=3,p=1$BCDndJ1kUOAAW/mwP7ViOQ$Ig4hpteBW1YM7Lrh3EHkHQ",
			"email1@company.com", "Pedro", "Petrenko", "77777777777", "false")

	mock.ExpectQuery(regexp.QuoteMeta(querySelectInfo)).
		WithArgs("user1").
		WillReturnRows(rowsInfo)

	user, err := userRepo.GetInfo("user1")
	assert.NoError(t, err)
	assert.Equal(t, &user1, user)

	rowsDisabled := sqlmock.NewRows(str).
		AddRow("3b60ac82-5e8f-4010-ac99-2344cfa72ce0", "user1",
			"$argon2id$v=19$m=65536,t=3,p=1$BCDndJ1kUOAAW/mwP7ViOQ$Ig4hpteBW1YM7Lrh3EHkHQ",
			"email1@company.com", "Pedro", "Petrenko", "77777777777", "true")

	mock.ExpectQuery(regexp.QuoteMeta(querySelectInfo)).
		WithArgs("user1").
		WillReturnRows(rowsDisabled)

	_, err = userRepo.GetInfo("user1")

	assert.Error(t, err)

}

func TestUpdatePassword(t *testing.T) {
	mock, userRepo, db := mockUserRepo(t)
	defer db.Close()
	timestamp := time.Now()
	mock.ExpectExec(regexp.QuoteMeta(queryUpdatePassword)).
		WithArgs("$argon2id$v=19$m=65536,t=3,p=1$RI2osB82TQY0w2gC3fitFQ$U9qrncg+AgvyIGwIeJZzmQ", &timestamp, "user1").
		WillReturnResult(driver.RowsAffected(1))

	err := userRepo.UpdatePassword("user1", "password2")
	_, err = db.Exec(queryUpdatePassword, "$argon2id$v=19$m=65536,t=3,p=1$RI2osB82TQY0w2gC3fitFQ$U9qrncg+AgvyIGwIeJZzmQ", timestamp, "user1")
	assert.NoError(t, err)
}

func TestGetEmail(t *testing.T) {
	mock, userRepo, db := mockUserRepo(t)
	defer db.Close()
	rowsInfo := sqlmock.NewRows([]string{"email"}).
		AddRow("email1@company.com")

	mock.ExpectQuery(regexp.QuoteMeta(queryGetEmail)).
		WithArgs("user1").
		WillReturnRows(rowsInfo)

	pwd, err := userRepo.GetEmail("user1")
	assert.NoError(t, err)
	assert.Equal(t, "email1@company.com", pwd)
}

func TestSetActivationCode(t *testing.T) {
	mock, _, db := mockUserRepo(t)
	defer db.Close()
	code := "RFhMcVpuRmNTdk9GY3NzSHJDZkVvRUlz"
	timestamp := time.Now()

	mock.ExpectExec(regexp.QuoteMeta(queryDisableCode)).
		WithArgs("user1").
		WillReturnResult(driver.RowsAffected(1))

	_, err := db.Exec(queryDisableCode, "user1")
	assert.NoError(t, err)
	mock.ExpectExec(regexp.QuoteMeta(querySetCode)).
		WithArgs("user1", code, timestamp.Add(time.Hour*24), true).
		WillReturnResult(driver.RowsAffected(1))

	_, err = db.Exec(querySetCode, "user1", code, timestamp.Add(time.Hour*24), true)

	assert.NoError(t, err)

}

func TestCheckActivationCode(t *testing.T) {
	mock, userRepo, db := mockUserRepo(t)
	defer db.Close()
	code := "RFhMcVpuRmNTdk9GY3NzSHJDZkVvRUlz"
	timestamp := time.Now()

	tests := []struct {
		name       string
		time       time.Time
		wantPValue bool
	}{
		{name: "good", time: timestamp.Add(time.Hour * 2), wantPValue: true},
		{name: "expired", time: timestamp.Add(time.Hour * (-2)), wantPValue: false},
	}

	for _, tt := range tests {
		rowsInfo := sqlmock.NewRows([]string{"code", "expired_at"}).
			AddRow(code, tt.time)

		mock.ExpectQuery(regexp.QuoteMeta(queryCheckCode)).
			WithArgs("user1").
			WillReturnRows(rowsInfo)

		active, err := userRepo.CheckActivationCode("user1", code)
		assert.NoError(t, err)
		assert.Equal(t, tt.wantPValue, active)
	}
}

func mockUserRepo(t *testing.T) (sqlmock.Sqlmock, *usersRepo, *sql.DB) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	return mock, SetUsersRepo(db), db
}

func mockUser() User {
	return User{
		ID:        "3b60ac82-5e8f-4010-ac99-2344cfa72ce0",
		Username:  "user1",
		Password:  "$argon2id$v=19$m=65536,t=3,p=1$BCDndJ1kUOAAW/mwP7ViOQ$Ig4hpteBW1YM7Lrh3EHkHQ",
		Email:     "email1@company.com",
		FirstName: "Pedro",
		LastName:  "Petrenko",
		Phone:     "77777777777",
	}
}
