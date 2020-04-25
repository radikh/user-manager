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
	mock, userRepo, db := mockUserRepo(t)
	defer db.Close()
	user := mockUser()

	mock.ExpectExec(regexp.QuoteMeta(queryInsert)).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := userRepo.Add(&user)
	assert.NoError(t, err)
}

func TestUpdate(t *testing.T) {
	mock, userRepo, db := mockUserRepo(t)
	defer db.Close()
	mock.ExpectExec(regexp.QuoteMeta(queryUpdate)).
		WillReturnResult(driver.RowsAffected(1))
	user := mockUser()
	err := userRepo.Update(&user)
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
	user1.Password = "$argon2id$v=19$m=65536,t=3,p=1$BCDndJ1kUOAAW/mwP7ViOQ$Ig4hpteBW1YM7Lrh3EHkHQ"
	user, err := userRepo.GetInfo("user1")
	assert.NoError(t, err)
	assert.Equal(t, &user1, user)

	mock.ExpectQuery(regexp.QuoteMeta(querySelectInfo)).
		WithArgs("user1").
		WillReturnError(sql.ErrNoRows)
	_, err = userRepo.GetInfo("user1")
	assert.EqualError(t, err, "There is no such user in database: sql: no rows in result set")

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
	mock.ExpectExec(regexp.QuoteMeta(queryUpdatePassword)).
		WillReturnResult(driver.RowsAffected(1))

	err := userRepo.UpdatePassword("user1", "password2")
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
	mock, userRepo, db := mockUserRepo(t)
	defer db.Close()
	code := "RFhMcVpuRmNTdk9GY3NzSHJDZkVvRUlz"
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(queryDisable)).
		WithArgs("true", "user1").
		WillReturnResult(driver.RowsAffected(1))

	mock.ExpectExec(regexp.QuoteMeta(queryDisableCode)).
		WithArgs("user1").
		WillReturnResult(driver.RowsAffected(1))

	mock.ExpectExec(regexp.QuoteMeta(querySetCode)).
		WillReturnResult(driver.RowsAffected(1))
	mock.ExpectCommit()
	err := userRepo.SetActivationCode("user1", code)

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

func mockUserRepo(t *testing.T) (sqlmock.Sqlmock, *UsersRepo, *sql.DB) {
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
		Password:  "password",
		Email:     "email1@company.com",
		FirstName: "Pedro",
		LastName:  "Petrenko",
		Phone:     "77777777777",
	}
}
