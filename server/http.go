// Package http presents HTTP server implementation.
// It provides a REST API to perform a set of CRUD
// to manage users and an endpoint to authenticate.
package server

import (
	"net/http"

	"github.com/lvl484/user-manager/logger"
	"github.com/lvl484/user-manager/middleware"
	"github.com/lvl484/user-manager/model"

	"github.com/gorilla/mux"
)

// HTTP server struct
type HTTP struct {
	address string
	acc     *account
	ur      *model.UsersRepo
}

// NewHTTP get address of server and return pointer to newserver
func NewHTTP(addr string, ur *model.UsersRepo) *HTTP {
	return &HTTP{
		address: addr,
		ur:      ur,
	}
}

// Start create all routes and starting server
func (h *HTTP) Start() error {
	mainRoute := mux.NewRouter()
	mainRoute.Use(middleware.NewBasicAuthentication(h.ur).Middleware)
	mainRoute.HandleFunc("/account", h.acc.CreateAccount).Methods("POST")
	mainRoute.HandleFunc("/account", h.acc.GetInfoAccount).Methods("GET")
	mainRoute.HandleFunc("/account", h.acc.UpdateAccount).Methods("PUT")
	mainRoute.HandleFunc("/account", h.acc.DeleteAccount).Methods("DELETE")
	mainRoute.HandleFunc("/validate", h.acc.ValidateAccount).Methods("GET")

	logger.LogUM.Infof("Server Listening at %s...", h.address)
	return http.ListenAndServe(h.address, mainRoute)
}
