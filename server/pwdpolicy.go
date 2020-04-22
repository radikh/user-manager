// Package http presents HTTP server implementation.
// It provides a REST API to perform a set of CRUD
// to manage users and an endpoint to authenticate.
package server

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"gopkg.in/gomail.v2"

	"github.com/lvl484/user-manager/logger"
	"github.com/lvl484/user-manager/model"
)

const (
	StatusCreateOK        = "Successfully created"
	StatusUpdateOK        = "Successfully updated"
	StatusRefreshOK       = "Successfully refreshed"
	StatusBadRequest      = "Bad request"
	StatusUnexpectedError = "Unexpected error"
	mailText              = `<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<title>User Manager password change</title>
	</head>
	<body>
		<h1>User Manager Service</h1>
		<h3>Thank you for requsting to change password!</h3>
	
		<p>Hello, You want to change password on User manager service.</p>
		<p>You need to do one more step to finish procedure...</p>
		<br>
		<p>Activation code: <strong>%s</strong></p>
		<p>
			Please note that activation code will be actual for 24 hours
		</p>
		<br>
		<p>If you did not request for changing password in User manager Service, please ignore this email!</p>
		<br>
		<p>
			<i>
				Best Regards,<br><br>
				User Manager Service<br>
			</i>
		</p>
	</body>
	</html>`
)

const (
	errNoUserName  = "No user name in header"
	errCodeNoMatch = "Code does not match"
	reqCode        = "User's <<%s>> make request for changing password"
	refCode        = "User's <<%s>> make request for refreshing activation code"
	updPassword    = "User's <<%s>> password successfully updated"
	runes          = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

// Police struct for changing password
type Policy struct {
	// Unique login
	Username string `json:"user_name"`
	// Strong enough password
	Password string `json:"password"`
	// Activation code for changing password
	Code string `json:"code"`
}

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

// decodePolicyFromBody draws up Policy structure from reguest body
func decodePolicyFromBody(r *http.Request) (*Policy, error) {
	var pwd *Policy
	err := json.NewDecoder(r.Body).Decode(&pwd)
	return pwd, err
}

// RequestPasswordChange start procedure for changing password:
// -disable accounnt for a time of passing of procedure
// -generate and send activation code to user's email
func (a *account) RequestPasswordChange(w http.ResponseWriter, r *http.Request) {
	login, err := a.createNewActivationCode(r)
	if err != nil {
		createErrorResponse(w, http.StatusBadRequest, StatusBadRequest, err)
		return
	}

	createJSONResponse(w, http.StatusOK, StatusCreateOK, fmt.Sprintf(reqCode, login))
}

// RefreshActivetionCode create new activation code and send to user
func (a *account) RefreshActivationCode(w http.ResponseWriter, r *http.Request) {
	login, err := a.createNewActivationCode(r)
	if err != nil {
		createErrorResponse(w, http.StatusBadRequest, StatusBadRequest, err)
		return
	}
	createJSONResponse(w, http.StatusOK, StatusCreateOK, fmt.Sprintf(refCode, login))
}

// createNewActivationCode create new activation code and send to user
func (a *account) createNewActivationCode(r *http.Request) (string, error) {
	login, _, ok := r.BasicAuth()
	if !ok {
		return "", errors.New(errNoUserName)
	}
	code := createCode(login)
	err := (*model.UsersRepo)(a).SetActivationCode(login, code)
	if err != nil {
		return "", err
	}
	email, err := (*model.UsersRepo)(a).GetEmail(login)
	if err != nil {
		return "", err
	}
	err = sendEmail(code, email)
	if err != nil {
		return "", err
	}
	return login, nil
}

// UpdatePassword ends procedure of changing password and update user's password
func (a *account) UpdatePassword(w http.ResponseWriter, r *http.Request) {
	var pwd *Policy

	pwd, err := decodePolicyFromBody(r)
	if err != nil {
		createErrorResponse(w, http.StatusBadRequest, StatusBadRequest, err)
		return
	}

	status, err := (*model.UsersRepo)(a).CheckActivationCode(pwd.Username, pwd.Code)
	if err != nil {
		createErrorResponse(w, http.StatusBadRequest, StatusBadRequest, err)
		return
	}
	if !status {
		createErrorResponse(w, http.StatusUnavailableForLegalReasons, StatusBadRequest, errors.New(errCodeNoMatch))
		return
	}
	err = (*model.UsersRepo)(a).UpdatePassword(pwd.Username, pwd.Password)
	if err != nil {
		createErrorResponse(w, http.StatusBadRequest, StatusBadRequest, err)
		return
	}

	createJSONResponse(w, http.StatusOK, StatusUpdateOK, fmt.Sprintf(updPassword, pwd.Username))
}

// sendEmail send activation code to user's email
func sendEmail(code string, mail string) error {
	email := gomail.NewMessage()

	email.SetHeader("From", "user.namager@gmail.com")
	email.SetHeader("To", mail)
	email.SetHeader("Subject", "Activation code for changing password")
	email.SetBody("text/html", fmt.Sprintf(mailText, code))

	dialer := gomail.NewDialer("smtp.gmail.com", 587, "user.namager@gmail.com", "lvl484Golang")

	return dialer.DialAndSend(email)
}

// createCode creates random activation code
func createCode(login string) string {
	str := fmt.Sprintf("%s%s", runes, login)
	var letters = []rune(str)

	rand.Seed(time.Now().UnixNano())
	b := make([]rune, 24)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	bs := []byte(string(b))
	return base64.RawStdEncoding.EncodeToString(bs)
}
