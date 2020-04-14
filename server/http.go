// Package http presents HTTP server implementation.
// It provides a REST API to perform a set of CRUD
// to manage users and an endpoint to authenticate.
package server

import (
	"context"
	"net/http"

	"github.com/lvl484/user-manager/config"
	"github.com/lvl484/user-manager/logger"
	"github.com/lvl484/user-manager/model"
	"github.com/lvl484/user-manager/server/http/middleware"

	"github.com/gorilla/mux"
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
func (h *HTTP) NewRouter() *mux.Router {
	mainRoute := mux.NewRouter()
	withoutAuth := mainRoute.PathPrefix("/account")
	withoutAuth.HandlerFunc(h.acc.CreateAccount).Methods(http.MethodPost)

	auth := mainRoute.PathPrefix("/um").Subrouter()
	auth.Use(middleware.NewBasicAuthentication(h.ur).Middleware)
	auth.HandleFunc("/account", h.acc.GetInfoAccount).Methods(http.MethodGet)
	auth.HandleFunc("/account", h.acc.UpdateAccount).Methods(http.MethodPut)
	auth.HandleFunc("/account", h.acc.DeleteAccount).Methods(http.MethodDelete)
	auth.HandleFunc("/validate", h.acc.ValidateAccount).Methods(http.MethodGet)

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
