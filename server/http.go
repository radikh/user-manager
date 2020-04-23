// Package http presents HTTP server implementation.
// It provides a REST API to perform a set of CRUD
// to manage users and an endpoint to authenticate.
package server

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/lvl484/user-manager/config"
	"github.com/lvl484/user-manager/logger"
	"github.com/lvl484/user-manager/model"
)

// HTTP server struct
type HTTP struct {
	srv *http.Server
	acc *account
	ur  *model.UsersRepo
}

// NewHTTP get address of server and return pointer to newserver
func NewHTTP(cfg *config.Config, ur *model.UsersRepo) *HTTP {
	srv := &http.Server{
		Addr:         cfg.ServerAddress(),
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	}

	return &HTTP{
		srv: srv,
		ur:  ur,
		acc: (*account)(ur),
	}
}

// NewRouter return new mux router
func (h *HTTP) NewRouter() *mux.Router {

	mainRoute := mux.NewRouter()

	mainRoute.HandleFunc("/password", h.acc.RequestPasswordChange).Methods(http.MethodPost)
	mainRoute.HandleFunc("/password", h.acc.UpdatePassword).Methods(http.MethodPut)
	mainRoute.HandleFunc("/password", h.acc.RefreshActivationCode).Methods(http.MethodGet)

	return mainRoute
}

// Start create all routes and starting server
func (h *HTTP) Start() error {
	logger.LogUM.Infof("Server Listening at %s...", h.srv.Addr)
	return http.ListenAndServe(h.srv.Addr, h.NewRouter())
}

// Stop stops all routes and stopping server
func (h *HTTP) Stop(ctx context.Context) error {
	return h.srv.Shutdown(ctx)
}
