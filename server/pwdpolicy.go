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

	"github.com/lvl484/user-manager/model"
)

const (
	errNoUserName       = "No user name in header"
	errCodeNoMatch      = "Code does not match"
	reqCode             = "User's ( %s ) make request for changing password"
	refCode             = "User's ( %s ) make request for refreshing activation code"
	updPassword         = "User's ( %s ) password successfully updated"
	loggerMessage       = `Operation responce code: "%d" Message: "%s" Data: "%v"`
	runes               = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	StatusRefreshOK     = "Successfully refreshed"
	StatusActionOK      = "Operation successfull"
	msgActiveCodeExists = "Active code exists. Check you's email or refresh for new"
	mailText            = `<!DOCTYPE html>
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

// Police struct for changing password
type Policy struct {
	// Unique login
	Username string `json:"user_name"`
	// Strong enough password
	Password string `json:"password"`
	// Activation code for changing password
	Code string `json:"code"`
}

// decodePolicyFromBody draws up Policy structure from reguest body
func decodePolicyFromBody(r *http.Request) (*Policy, error) {
	var pwd *Policy
	err := json.NewDecoder(r.Body).Decode(&pwd)
	return pwd, err
}

// decodeLoginFromHeader draws up username from header
func decodeLoginFromHeader(r *http.Request) (string, error) {
	login, _, ok := r.BasicAuth()
	if !ok {
		return "", errors.New(errNoUserName)
	}
	return login, nil
}

// RequestPasswordChange start procedure for changing password:
// -disable accounnt for a time of passing of procedure
// -generate and send activation code to user's email
func (a *account) RequestPasswordChange(w http.ResponseWriter, r *http.Request) {
	login, err := decodeLoginFromHeader(r)
	if err != nil {
		createErrorResponse(w, http.StatusBadRequest, StatusBadRequest, err)
		return
	}
	_, err = (*model.UsersRepo)(a).CheckCodeForUser(login)
	if err == nil {
		createErrorResponse(w, http.StatusBadRequest, StatusBadRequest, errors.New(msgActiveCodeExists))
		return
	}
	err = a.createNewActivationCode(login)
	if err != nil {
		createErrorResponse(w, http.StatusBadRequest, StatusBadRequest, err)
		return
	}

	createJSONResponse(w, http.StatusOK, StatusActionOK, fmt.Sprintf(reqCode, login))
}

// RefreshActivetionCode create new activation code and send to user
func (a *account) RefreshActivationCode(w http.ResponseWriter, r *http.Request) {
	login, err := decodeLoginFromHeader(r)
	if err != nil {
		createErrorResponse(w, http.StatusBadRequest, StatusBadRequest, err)
		return
	}
	err = a.createNewActivationCode(login)
	if err != nil {
		createErrorResponse(w, http.StatusBadRequest, StatusBadRequest, err)
		return
	}
	createJSONResponse(w, http.StatusOK, StatusActionOK, fmt.Sprintf(refCode, login))
}

// createNewActivationCode create new activation code and send to user
func (a *account) createNewActivationCode(login string) error {
	code := createCode(login)
	err := (*model.UsersRepo)(a).SetActivationCode(login, code)
	if err != nil {
		return err
	}
	email, err := (*model.UsersRepo)(a).GetEmail(login)
	if err != nil {
		return err
	}
	err = sendEmail(code, email)
	if err != nil {
		return err
	}
	return nil
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
	// TODO configuration data as sender, host, password would be taken from conf module after it updates
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
