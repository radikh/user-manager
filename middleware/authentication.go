// Package middleware provides different implementations of HTTP middleware.
// Each implementation should be of type MiddlewareFunc
// https://pkg.go.dev/github.com/gorilla/mux?tab=doc#MiddlewareFunc
//
// Example
//
// A very basic middleware which logs the URI of the request being handled could be written as
//
//		func simpleMw(next http.Handler) http.Handler {
//			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//				// Do stuff here
//				log.Println(r.RequestURI)
//				// Call the next handler, which can be another middleware in the chain, or the final handler.
//				next.ServeHTTP(w, r)
//			})
//		}
//
// Middleware can be added to a router using `Router.Use()`
//
//		r := mux.NewRouter()
//		r.HandleFunc("/", handler)
//		r.Use(simpleMw)
//
// Source: https://pkg.go.dev/github.com/gorilla/mux?tab=doc#pkg-overview
package middleware

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/lvl484/user-manager/logger"
	"github.com/lvl484/user-manager/model"
)

const messageUnauthorized = "Authenticate failed"

type UserProvider interface {
	GetInfo(username string) (*model.User, error)
}

type BasicAuthentication struct {
	ur UserProvider
}

func NewBasicAuthentication(ur UserProvider) *BasicAuthentication {
	return &BasicAuthentication{ur: ur}
}

func (a *BasicAuthentication) Middleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()
		if !ok {
			unauthorized(w)
			return
		}

		userFromDB, err := a.ur.GetInfo(user)
		if err != nil {
			internalServerError(w, err)
			return
		}

		matched, err := model.ComparePassword(pass, userFromDB.Password)
		if err != nil {
			internalServerError(w, err)
			return
		}

		if !matched {
			unauthorized(w)
			return
		}

		logger.LogUM.Debugf("Authentication successful! Hello, %s", user)

		handler.ServeHTTP(w, r)
	})
}

func unauthorized(w http.ResponseWriter) {
	w.Header().Set("WWW-Authenticate", `Basic realm="user-manager"`)
	w.WriteHeader(http.StatusUnauthorized)

	authError := &model.Error{
		Code:    strconv.Itoa(http.StatusUnauthorized),
		Message: messageUnauthorized,
	}

	err := json.NewEncoder(w).Encode(&authError)
	if err != nil {
		logger.LogUM.Errorf("Write unauthorized response error: %v", err)
	}

	logger.LogUM.Info("Authentication failed! Invalid login or password")
}

func internalServerError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)

	internalError := &model.Error{
		Code:    strconv.Itoa(http.StatusInternalServerError),
		Message: err.Error(),
	}

	err = json.NewEncoder(w).Encode(&internalError)
	if err != nil {
		logger.LogUM.Errorf("Write internal server response error: %v", err)
	}

	logger.LogUM.Error("Internal server error")
}
