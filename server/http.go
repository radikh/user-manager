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

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type HTTP struct {
	srv *http.Server
	ur  *model.UsersRepo
}

func NewHTTP(cfg *config.Config, ur *model.UsersRepo) *HTTP {
	srv := &http.Server{
		Addr:         cfg.ServerAddress(),
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	}

	return &HTTP{
		srv: srv,
		ur:  ur,
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

	h.srv.Handler = mainRoute

	logger.LogUM.Infof("Server Listening at %s...", h.srv.Addr)
	return h.srv.ListenAndServe()
}

// Stop stops all routes and stopping server
func (h *HTTP) Stop(ctx context.Context) error {
	return h.srv.Shutdown(ctx)
}
