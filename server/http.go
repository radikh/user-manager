// Package http presents HTTP server implementation.
// It provides a REST API to perform a set of CRUD
// to manage users and an endpoint to authenticate.
package server

import (
	"net/http"

	"github.com/lvl484/user-manager/logger"
	"github.com/lvl484/user-manager/middleware"
	"github.com/lvl484/user-manager/model"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type HTTP struct {
	address string
	ur      *model.UsersRepo
}

func NewHTTP(addr string, ur *model.UsersRepo) *HTTP {
	return &HTTP{
		address: addr,
		ur:      ur,
	}
}

// TODO: this method is using only for testing work
func (h *HTTP) UUID(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte(uuid.New().String()))
}

// Start create all routes and starting server
func (h *HTTP) Start() error {
	mainRoute := mux.NewRouter()
	mainRoute.Use(middleware.NewBasicAuthentication(h.ur).Middleware)
	// TODO: replace it with necessary REST APIs
	mainRoute.HandleFunc("/uuid", h.UUID).Methods(http.MethodGet)

	logger.LogUM.Infof("Server Listening at %s...", h.address)
	return http.ListenAndServe(h.address, mainRoute)
}
