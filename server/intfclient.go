// Package http presents HTTP server implementation.
// It provides a REST API to perform a set of CRUD
// to manage users and an endpoint to authenticate.
package server

import (
	"net/http"
)

type Account interface {
	CreateAccount(w http.ResponseWriter, r *http.Request)
	GetInfoAccount(w http.ResponseWriter, r *http.Request)
	UpdateAccount(w http.ResponseWriter, r *http.Request)
	DeleteAccount(w http.ResponseWriter, r *http.Request)
	ValidateAccount(w http.ResponseWriter, r *http.Request)
}
