// Command umcli provides admin command line tool to manipulate accounts with admin rights.
package pgclient

import (
	"encoding/json"
	"net/http"

	"github.com/lvl484/user-manager/model"
)

type Account interface {
	CreateAccount(w http.ResponseWriter, r *http.Request)
	GetInfoAccount(w http.ResponseWriter, r *http.Request)
	UpdateAccount(w http.ResponseWriter, r *http.Request)
	DeleteAccount(w http.ResponseWriter, r *http.Request)
}

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
	w.Write([]byte(StatusCreateOK.message))
}

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
	w.Write([]byte(StatusInfoOK.message))
}

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
	w.Write([]byte(StatusUpdateOK.message))
}

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
	w.Write([]byte(StatusDeleteOK.message))
}
