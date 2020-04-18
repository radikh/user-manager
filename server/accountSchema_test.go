package server

import (
	"testing"

	"github.com/lvl484/user-manager/model"
	"github.com/stretchr/testify/assert"
)

func Test_convertToResponsCreateAccount(t *testing.T) {
	user := model.User{
		ID:        "6e10da63-7c99-43be-8f05-3c7e8583ef0d",
		Username:  "nick",
		Password:  "pass",
		Email:     "n@sd.d",
		FirstName: "Mykola",
		LastName:  "Mykolajchuk",
		Phone:     "0001112233",
		CreatedAt: nil,
		UpdatedAt: nil,
	}
	createAccount := responseCreateAccount{
		ID:        "6e10da63-7c99-43be-8f05-3c7e8583ef0d",
		Username:  "nick",
		Email:     "n@sd.d",
		FirstName: "Mykola",
		LastName:  "Mykolajchuk",
		Phone:     "0001112233",
	}
	userZero := model.User{}
	createAccountZero := responseCreateAccount{}
	get := convertToResponseCreateAccount(&user)
	assert.Equal(t, get, createAccount)
	assert.NotEqual(t, get, createAccountZero)
	assert.NotNil(t, get)

	get = convertToResponseCreateAccount(&userZero)
	assert.Equal(t, get, createAccountZero)
	assert.NotEqual(t, get, createAccount)
	assert.NotNil(t, get)
}

func Test_convertToResponseAccountInfo(t *testing.T) {
	user := model.User{
		ID:        "6e10da63-7c99-43be-8f05-3c7e8583ef0d",
		Username:  "nick",
		Password:  "pass",
		Email:     "n@sd.d",
		FirstName: "Mykola",
		LastName:  "Mykolajchuk",
		Phone:     "0001112233",
		CreatedAt: nil,
		UpdatedAt: nil,
	}
	infoAccount := responseAccountInfo{
		Username:  "nick",
		Email:     "n@sd.d",
		FirstName: "Mykola",
		LastName:  "Mykolajchuk",
		Phone:     "0001112233",
	}
	userZero := model.User{}
	infoAccountZero := responseAccountInfo{}

	get := convertToResponseAccountInfo(&user)
	assert.Equal(t, get, infoAccount)
	assert.NotEqual(t, get, infoAccountZero)
	assert.NotNil(t, get)

	get = convertToResponseAccountInfo(&userZero)
	assert.Equal(t, get, infoAccountZero)
	assert.NotEqual(t, get, infoAccount)
	assert.NotNil(t, get)
}
