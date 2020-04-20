// Package http presents HTTP server implementation.
// It provides a REST API to perform a set of CRUD
// to manage users and an endpoint to authenticate.
package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"time"

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

	StatusVerificationOK  = "Successfully verified"
	VerificationLiveHours = 24

	htmlVerification     = "verification.html"
	htmlVerificationPath = "server/mail/mail_template/verification.html"
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

	createJSONResponse(w, code, msg, nil)
}

// decodeUserFromBody draws up user structure from request body
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
		return
	}

	err = (*model.UsersRepo)(a).Add(user)
	if err != nil {
		createErrorResponse(w, http.StatusBadRequest, StatusBadRequest, err)
		return
	}

	createJSONResponse(w, http.StatusCreated, StatusCreateOK, user)
}

// GetInfoAccount check if account exist and return info about user
func (a *account) GetInfoAccount(w http.ResponseWriter, r *http.Request) {
	user, err := decodeUserFromBody(w, r)
	if err != nil {
		return
	}

	user, err = (*model.UsersRepo)(a).GetInfo(user.Username)
	if err != nil {
		createErrorResponse(w, http.StatusBadRequest, StatusBadRequest, err)
		return
	}

	createJSONResponse(w, http.StatusOK, StatusInfoOK, user)
}

// UpdateAccount update data of account
func (a *account) UpdateAccount(w http.ResponseWriter, r *http.Request) {
	user, err := decodeUserFromBody(w, r)
	if err != nil {
		return
	}

	err = (*model.UsersRepo)(a).Update(user)
	if err != nil {
		createErrorResponse(w, http.StatusBadRequest, StatusBadRequest, err)
		return
	}

	createJSONResponse(w, http.StatusOK, StatusUpdateOK, user)
}

// DeleteAccount deletes account of user in database
func (a *account) DeleteAccount(w http.ResponseWriter, r *http.Request) {
	user, err := decodeUserFromBody(w, r)
	if err != nil {
		return
	}

	err = (*model.UsersRepo)(a).Delete(user.Username)
	if err != nil {
		createErrorResponse(w, http.StatusBadRequest, StatusBadRequest, err)
		return
	}

	createJSONResponse(w, http.StatusNoContent, StatusDeleteOK, nil)
}

// ValidateAccount check if such account exist, check password and return user's info
func (a *account) ValidateAccount(w http.ResponseWriter, r *http.Request) {
	user, err := decodeUserFromBody(w, r)
	if err != nil {
		return
	}

	dbuser, err := (*model.UsersRepo)(a).GetInfo(user.Username)
	if err != nil {
		createErrorResponse(w, http.StatusBadRequest, StatusBadRequest, err)
		return
	}

	pwdValid, err := model.ComparePassword(user.Password, dbuser.Password)
	if err != nil {
		createErrorResponse(w, http.StatusBadRequest, StatusBadRequest, err)
		return
	}

	if !pwdValid {
		createErrorResponse(w, http.StatusUnauthorized, StatusAuthenticateFailed, nil)
		return
	}

	createJSONResponse(w, http.StatusOK, StatusInfoOK, dbuser)
}

// VerificationAccount gets verification info from request body (login, code and password)
func (a *account) VerificationAccount(w http.ResponseWriter, r *http.Request) {
	verification, err := decodeVerificationCodeFromBody(w, r)
	if err != nil {
		createErrorResponse(w, http.StatusBadRequest, StatusBadRequest, err)
		return
	}

	vt, vc, err := (*model.UsersRepo)(a).GetVerificationCodeTime(verification.Login)
	if err != nil {
		createErrorResponse(w, http.StatusBadRequest, StatusBadRequest, err)
		return
	}

	if time.Now().Hour()-vt.Hour() > VerificationLiveHours {
		createErrorResponse(w, http.StatusBadRequest, StatusBadRequest, fmt.Errorf("verification code expired"))
		return
	}

	if verification.Code == vc {
		dbuser, err := (*model.UsersRepo)(a).GetUserInfoIncludingSalted(verification.Login)
		if err != nil {
			createErrorResponse(w, http.StatusBadRequest, StatusBadRequest, err)
			return
		}

		pwdValid, err := model.ComparePassword(verification.Password, dbuser.Password)
		if err != nil {
			createErrorResponse(w, http.StatusBadRequest, StatusBadRequest, fmt.Errorf("verification password mismatched"))
			return
		}

		if !pwdValid {
			createErrorResponse(w, http.StatusUnauthorized, StatusAuthenticateFailed, fmt.Errorf("wrong password"))
			return
		}

		err = (*model.UsersRepo)(a).Activate(verification.Login)
		if err != nil {
			createErrorResponse(w, http.StatusBadRequest, StatusBadRequest, err)
			return
		}

		createJSONResponse(w, http.StatusOK, StatusVerificationOK, verification)
		return
	}

	createErrorResponse(w, http.StatusBadRequest, StatusBadRequest, fmt.Errorf("verification code is invalid"))
}

// decodeVerificationCodeFromBody fills verification structure from request body
func decodeVerificationCodeFromBody(w http.ResponseWriter, r *http.Request) (*model.Verification, error) {
	var verification *model.Verification

	err := json.NewDecoder(r.Body).Decode(&verification)
	if err != nil {
		createErrorResponse(w, http.StatusBadRequest, StatusBadRequest, err)
	}

	return verification, err
}

// ReadHTMLVerificationPage reads HTML page from template and writes it into response
func (a *account) ReadHTMLVerificationPage(w http.ResponseWriter, r *http.Request) {
	page := template.New(htmlVerification)

	page, err := page.ParseFiles(htmlVerificationPath)
	if err != nil {
		createErrorResponse(w, http.StatusInternalServerError, StatusUnexpectedError, err)
	}

	data := struct {
		Login string
		Code  string
	}{
		Login: r.URL.Query().Get("login"),
		Code:  r.URL.Query().Get("code"),
	}

	var tpl bytes.Buffer
	if err := page.Execute(&tpl, data); err != nil {
		createErrorResponse(w, http.StatusInternalServerError, StatusUnexpectedError, err)
	}

	_, err = w.Write(tpl.Bytes())
	if err != nil {
		createErrorResponse(w, http.StatusInternalServerError, StatusUnexpectedError, err)
	}
}
