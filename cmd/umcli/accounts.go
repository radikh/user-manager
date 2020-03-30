// Command umcli provides admin command line tool to manipulate accounts with admin rights.
package main

import (
	"encoding/json"
	"net/http"

	"github.com/lvl484/user-manager/logger"
	"github.com/lvl484/user-manager/model"
)

// StatusAccount structure that hold responce of User Manager client subservice
type StatusAccount struct {
	code    int
	message string
}

// Response's
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

// CreateAccount create a new account in database
func (ur *usersRepo) CreateAccount(w http.ResponseWriter, r *http.Request) {
	var user *model.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, StatusBadRequest.message, StatusBadRequest.code)
		return
	}
	loginExist, err := ur.CheckLoginExist(user.Username)
	if err != nil {
		http.Error(w, StatusUnexpectedError.message, StatusUnexpectedError.code)
		return
	}
	if loginExist {
		http.Error(w, StatusLoginInUse.message, StatusLoginInUse.code)
		return
	}
	err = ur.Add(user)
	if err != nil {
		http.Error(w, StatusUnexpectedError.message, StatusUnexpectedError.code)
		return
	}
	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		http.Error(w, StatusUnexpectedError.message, StatusUnexpectedError.code)
		return
	}
	w.WriteHeader(StatusCreateOK.code)
	_, err = w.Write([]byte(StatusCreateOK.message))
	if err != nil {
		logger.LogUM.Error(err)
	}
}

// GetInfoAccount check if account exist and return info about user
func (ur *usersRepo) GetInfoAccount(w http.ResponseWriter, r *http.Request) {
	var user *model.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, StatusBadRequest.message, StatusBadRequest.code)
		return
	}
	loginExist, err := ur.CheckLoginExist(user.Username)
	if err != nil {
		http.Error(w, StatusUnexpectedError.message, StatusUnexpectedError.code)
		return
	}
	if !loginExist {
		http.Error(w, StatusAccountNotExist.message, StatusAccountNotExist.code)
		return
	}
	user, err = ur.GetInfo(user.Username)
	if err != nil {
		http.Error(w, StatusUnexpectedError.message, StatusUnexpectedError.code)
		return
	}
	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		http.Error(w, StatusUnexpectedError.message, StatusUnexpectedError.code)
		return
	}
	w.WriteHeader(StatusInfoOK.code)
	_, err = w.Write([]byte(StatusInfoOK.message))
	if err != nil {
		logger.LogUM.Error(err)
	}
}

// UpdateAccount update data of account
func (ur *usersRepo) UpdateAccount(w http.ResponseWriter, r *http.Request) {
	var user *model.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, StatusBadRequest.message, StatusBadRequest.code)
		return
	}
	loginExist, err := ur.CheckLoginExist(user.Username)
	if err != nil {
		http.Error(w, StatusUnexpectedError.message, StatusUnexpectedError.code)
		return
	}
	if !loginExist {
		http.Error(w, StatusAuthenticateFailed.message, StatusAuthenticateFailed.code)
		return
	}
	err = ur.Update(user)
	if err != nil {
		http.Error(w, StatusUnexpectedError.message, StatusUnexpectedError.code)
		return
	}
	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		http.Error(w, StatusUnexpectedError.message, StatusUnexpectedError.code)
		return
	}
	w.WriteHeader(StatusUpdateOK.code)
	_, err = w.Write([]byte(StatusUpdateOK.message))
	if err != nil {
		logger.LogUM.Error(err)
	}
}

// DeleteAccount deletes account of user in database
func (ur *usersRepo) DeleteAccount(w http.ResponseWriter, r *http.Request) {
	var user *model.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, StatusBadRequest.message, StatusBadRequest.code)
		return
	}
	loginExist, err := ur.CheckLoginExist(user.Username)
	if err != nil {
		http.Error(w, StatusUnexpectedError.message, StatusUnexpectedError.code)
		return
	}
	if !loginExist {
		http.Error(w, StatusAccountNotExist.message, StatusAccountNotExist.code)
		return
	}
	err = ur.Delete(user.Username)
	if err != nil {
		http.Error(w, StatusUnexpectedError.message, StatusUnexpectedError.code)
		return
	}
	w.WriteHeader(StatusDeleteOK.code)
	_, err = w.Write([]byte(StatusDeleteOK.message))
	if err != nil {
		logger.LogUM.Error(err)
	}
}

// ValidateAccount check if such account exist, check password and return user's info
func (ur *usersRepo) ValidateAccount(w http.ResponseWriter, r *http.Request) {
	var user, dbuser *model.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, StatusBadRequest.message, StatusBadRequest.code)
		return
	}
	loginExist, err := ur.CheckLoginExist(user.Username)
	if err != nil {
		http.Error(w, StatusUnexpectedError.message, StatusUnexpectedError.code)
		return
	}
	if !loginExist {
		http.Error(w, StatusAccountNotExist.message, StatusAccountNotExist.code)
		return
	}

	user, err = ur.GetInfo(user.Username)
	if err != nil {
		http.Error(w, StatusUnexpectedError.message, StatusUnexpectedError.code)
		return
	}

	pwdValid, err := ComparePassword(user.Password, dbuser.Password)
	if err != nil {
		http.Error(w, StatusUnexpectedError.message, StatusUnexpectedError.code)
		return
	}
	if !pwdValid {
		http.Error(w, StatusAuthenticateFailed.message, StatusAuthenticateFailed.code)
		return
	}
	err = json.NewEncoder(w).Encode(dbuser)
	if err != nil {
		http.Error(w, StatusUnexpectedError.message, StatusUnexpectedError.code)
		return
	}
	w.WriteHeader(StatusInfoOK.code)
	_, err = w.Write([]byte(StatusInfoOK.message))
	if err != nil {
		logger.LogUM.Error(err)
	}
}
