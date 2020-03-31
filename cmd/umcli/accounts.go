// Command umcli provides admin command line tool to manipulate accounts with admin rights.
package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/lvl484/user-manager/logger"
	"github.com/lvl484/user-manager/model"
)

// StatusAccount structure that hold responce of User Manager client subservice
type StatusAccount struct {
	code    int
	message string
}

var (
	StatusInfoOK             = StatusAccount{code: 200, message: "An account info"}
	StatusCreateOK           = StatusAccount{code: 201, message: "Successfully created"}
	StatusUpdateOK           = StatusAccount{code: 200, message: "Successfully updated"}
	StatusDeleteOK           = StatusAccount{code: 204, message: "Successfully deleted"}
	StatusBadRequest         = StatusAccount{code: 400, message: "Bad request"}
	StatusAuthenticateFailed = StatusAccount{code: 401, message: "Authenticate failed"}
	StatusAccountNotExist    = StatusAccount{code: 404, message: "Account does not exist"}
	StatusLoginInUse         = StatusAccount{code: 409, message: "Login in use"}
	StatusUnexpectedError    = StatusAccount{code: 444, message: "Unexpected error"}
)

// JSON create a JSON responce
func JSON(w http.ResponseWriter, sa StatusAccount, data interface{}) {
	w.WriteHeader(sa.code)
	fmt.Fprintf(w, "%s", sa.message)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		logger.LogUM.Error(err)
	}
}

// ERROR create an error responce
func ERROR(w http.ResponseWriter, sa StatusAccount, err error) {
	if err != nil {
		JSON(w, sa, struct {
			Error string `json:"error"`
		}{
			Error: err.Error(),
		})
		return
	}
	JSON(w, sa, nil)
}

// CreateAccount create a new account in database
func (ur *usersRepo) CreateAccount(w http.ResponseWriter, r *http.Request) {
	var user *model.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		ERROR(w, StatusBadRequest, err)
		return
	}
	loginExist, err := ur.CheckLoginExist(user.Username)
	if err != nil {
		ERROR(w, StatusUnexpectedError, err)
		return
	}
	if loginExist {
		JSON(w, StatusLoginInUse, nil)
		return
	}
	err = ur.Add(user)
	if err != nil {
		ERROR(w, StatusUnexpectedError, err)
		return
	}
	JSON(w, StatusCreateOK, user)
}

// GetInfoAccount check if account exist and return info about user
func (ur *usersRepo) GetInfoAccount(w http.ResponseWriter, r *http.Request) {
	var user *model.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		ERROR(w, StatusBadRequest, err)
		return
	}
	loginExist, err := ur.CheckLoginExist(user.Username)
	if err != nil {
		ERROR(w, StatusUnexpectedError, err)
		return
	}
	if !loginExist {
		JSON(w, StatusAccountNotExist, nil)
		return
	}
	user, err = ur.GetInfo(user.Username)
	if err != nil {
		ERROR(w, StatusUnexpectedError, err)
		return
	}
	JSON(w, StatusInfoOK, user)
}

// UpdateAccount update data of account
func (ur *usersRepo) UpdateAccount(w http.ResponseWriter, r *http.Request) {
	var user *model.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		ERROR(w, StatusBadRequest, err)
		return
	}
	loginExist, err := ur.CheckLoginExist(user.Username)
	if err != nil {
		ERROR(w, StatusUnexpectedError, err)
		return
	}
	if !loginExist {
		JSON(w, StatusAccountNotExist, nil)
		return
	}
	err = ur.Update(user)
	if err != nil {
		ERROR(w, StatusUnexpectedError, err)
		return
	}
	JSON(w, StatusUpdateOK, user)
}

// DeleteAccount deletes account of user in database
func (ur *usersRepo) DeleteAccount(w http.ResponseWriter, r *http.Request) {
	var user *model.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		ERROR(w, StatusBadRequest, err)
		return
	}
	loginExist, err := ur.CheckLoginExist(user.Username)
	if err != nil {
		ERROR(w, StatusUnexpectedError, err)
		return
	}
	if !loginExist {
		JSON(w, StatusAccountNotExist, nil)
		return
	}
	err = ur.Delete(user.Username)
	if err != nil {
		ERROR(w, StatusUnexpectedError, err)
		return
	}
	JSON(w, StatusDeleteOK, nil)
}

// ValidateAccount check if such account exist, check password and return user's info
func (ur *usersRepo) ValidateAccount(w http.ResponseWriter, r *http.Request) {
	var user, dbuser *model.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		ERROR(w, StatusBadRequest, err)
		return
	}
	loginExist, err := ur.CheckLoginExist(user.Username)
	if err != nil {
		ERROR(w, StatusUnexpectedError, err)
		return
	}
	if !loginExist {
		JSON(w, StatusAccountNotExist, nil)
		return
	}

	user, err = ur.GetInfo(user.Username)
	if err != nil {
		ERROR(w, StatusUnexpectedError, err)
		return
	}

	pwdValid, err := ComparePassword(user.Password, dbuser.Password)
	if err != nil {
		ERROR(w, StatusUnexpectedError, err)
		return
	}
	if !pwdValid {
		JSON(w, StatusAuthenticateFailed, nil)
		return
	}
	JSON(w, StatusInfoOK, dbuser)

}
