// Package http presents HTTP server implementation.
// It provides a REST API to perform a set of CRUD
// to manage users and an endpoint to authenticate.
package server

import (
	"encoding/json"
	"net/http"

	"github.com/lvl484/user-manager/logger"
	"github.com/lvl484/user-manager/model"
)

const (
	StatusInfoOK             = "An account info"
	StatusCreateOK           = "Successfully created"
	StatusUpdateOK           = "Successfully updated"
	StatusDeleteOK           = "Successfully deleted"
	StatusBadRequest         = "Bad request"
	StatusAuthenticateFailed = "Authenticate failed"
	StatusAccountNotExist    = "Account does not exist"
	StatusLoginInUse         = "Login in use"
	StatusUnexpectedError    = "Unexpected error"
)

type account model.UsersRepo

// createJSONResponse create a JSON response
func createJSONResponse(w http.ResponseWriter, code int, msg string, data interface{}) {
	w.WriteHeader(code)
	_, err := w.Write([]byte(msg))

	if err != nil {
		logger.LogUM.Error(err)
	}

	err = json.NewEncoder(w).Encode(data)
	if err != nil {
		logger.LogUM.Error(err)
	}
}

// createErrorResponse create an error response
func createErrorResponse(w http.ResponseWriter, code int, msg string, err error) {
	if err != nil {
		createJSONResponse(w, code, msg, struct {
			Error string `json:"error"`
		}{
			Error: err.Error(),
		})

		return
	}

	createJSONResponse(w, code, msg, err)
}

// decodeUserFromBody draws up user structure from reguest body
func decodeUserFromBody(w http.ResponseWriter, r *http.Request) (*model.User, error) {
	var user *model.User
	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		createErrorResponse(w, http.StatusBadRequest, StatusBadRequest, err)
	}

	return user, err
}

// CreateAccount create a new account in database
func (a *account) CreateAccount(w http.ResponseWriter, r *http.Request) {
	user, err := decodeUserFromBody(w, r)
	if err != nil {
		createErrorResponse(w, http.StatusBadRequest, StatusBadRequest, err)
		return
	}

	err = (*model.UsersRepo)(a).Add(user)
	if err != nil {
		createErrorResponse(w, http.StatusBadRequest, StatusBadRequest, err)
		return
	}

	createJSONResponse(w, http.StatusCreated, StatusCreateOK, convertToResponseCreateAccount(user))
}

// GetInfoAccount check if account exist and return info about user
func (a *account) GetInfoAccount(w http.ResponseWriter, r *http.Request) {
	username, _, ok := r.BasicAuth()
	if !ok {
		createErrorResponse(w, http.StatusBadRequest, StatusBadRequest, nil)
		return
	}

	user, err := (*model.UsersRepo)(a).GetInfo(username)
	if err != nil {
		createErrorResponse(w, http.StatusBadRequest, StatusBadRequest, err)
		return
	}

	createJSONResponse(w, http.StatusOK, StatusInfoOK, convertToResponseAccountInfo(user))
}

// UpdateAccount update data of account
func (a *account) UpdateAccount(w http.ResponseWriter, r *http.Request) {
	var user *model.User

	username, password, ok := r.BasicAuth()
	if !ok {
		createErrorResponse(w, http.StatusBadRequest, StatusBadRequest, nil)
		return
	}

	user, err := decodeUserFromBody(w, r)
	if err != nil {
		createErrorResponse(w, http.StatusBadRequest, StatusBadRequest, nil)
		return
	}

	user.Username = username
	user.Password = password

	err = (*model.UsersRepo)(a).Update(user)
	if err != nil {
		createErrorResponse(w, http.StatusBadRequest, StatusBadRequest, err)
		return
	}

	createJSONResponse(w, http.StatusOK, StatusUpdateOK, convertToResponseAccountInfo(user))
}

// DeleteAccount deletes account of user in database
func (a *account) DeleteAccount(w http.ResponseWriter, r *http.Request) {
	username, _, ok := r.BasicAuth()
	if !ok {
		createErrorResponse(w, http.StatusBadRequest, StatusBadRequest, nil)
		return
	}

	err := (*model.UsersRepo)(a).Delete(username)
	if err != nil {
		createErrorResponse(w, http.StatusBadRequest, StatusBadRequest, err)
		return
	}

	createJSONResponse(w, http.StatusNoContent, StatusDeleteOK, nil)
}

// ValidateAccount check if such account exist, check password and return user ID
func (a *account) ValidateAccount(w http.ResponseWriter, r *http.Request) {
	username, password, ok := r.BasicAuth()
	if !ok {
		createErrorResponse(w, http.StatusBadRequest, StatusBadRequest, nil)
		return
	}

	user, err := (*model.UsersRepo)(a).GetInfo(username)

	if err != nil {
		createErrorResponse(w, http.StatusBadRequest, StatusBadRequest, err)
		return
	}

	pwdValid, err := model.ComparePassword(password, user.Password)
	if err != nil {
		createErrorResponse(w, http.StatusBadRequest, StatusBadRequest, err)
		return
	}

	if !pwdValid {
		createErrorResponse(w, http.StatusUnauthorized, StatusAuthenticateFailed, nil)
		return
	}

	createJSONResponse(w, http.StatusOK, StatusInfoOK, user.ID)
}
