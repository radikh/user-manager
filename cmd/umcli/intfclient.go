// Command umcli provides admin command line tool to manipulate accounts with admin rights.
package main

import (
	"net/http"

	"github.com/lvl484/user-manager/model"
)

type Account interface {
	CreateAccount(w http.ResponseWriter, r *http.Request)
	GetInfoAccount(w http.ResponseWriter, r *http.Request)
	UpdateAccount(w http.ResponseWriter, r *http.Request)
	DeleteAccount(w http.ResponseWriter, r *http.Request)
	ValidateAccount(w http.ResponseWriter, r *http.Request)
}

type Users interface {
	Add(args ...interface{}) error
	Update(args ...interface{}) error
	Delete(login string) error
	Disable(login string) error
	Activate(login string) error
	GetInfo(login string) (*model.User, error)
	CheckLoginExist(lo string) (bool, error)
}
