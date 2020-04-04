// Package http presents HTTP server implementation.
// It provides a REST API to perform a set of CRUD
// to manage users and an endpoint to authenticate.
package server

import (
	"github.com/gorilla/mux"
)

func (account *Account) NewRoute() *mux.Router {
	mainRoute := mux.NewRouter()
	//mainRoute.Use(middleware)
	apiRoute := mainRoute.PathPrefix("um").Subrouter()
	apiRoute.HandleFunc("/account", account.CreateAccount).Methods("POST")
	apiRoute.HandleFunc("/account", account.GetInfoAccount).Methods("GET")
	apiRoute.HandleFunc("/account", account.UpdateAccount).Methods("PUT")
	apiRoute.HandleFunc("/account", account.DeleteAccount).Methods("DELETE")
	apiRoute.HandleFunc("/validate", account.ValidateAccount).Methods("GET")

	return mainRoute
}
