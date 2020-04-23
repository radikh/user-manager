// Package http presents HTTP server implementation.
// It provides a REST API to perform a set of CRUD
// to manage users and an endpoint to authenticate.
package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/lvl484/user-manager/logger"
	"github.com/lvl484/user-manager/model"
)

const (
	StatusInfoOK             = "An account info"
	StatusCreateOK           = "Successfully created"
	StatusUpdateOK           = "Successfully updated"
	StatusDeleteOK           = "Successfully deleted"
	StatusBadRequest         = "Bad request"
	StatusAuthenticateFailed = "Authenticate failed"
	StatusAccountNotExist    = "Account does not exist"
	StatusLoginInUse         = "Login in use"
	StatusUnexpectedError    = "Unexpected error"
)

type account model.UsersRepo

// createJSONResponse create a JSON response
func createJSONResponse(w http.ResponseWriter, code int, msg string, data interface{}) {
	w.WriteHeader(code)
	_, err := w.Write([]byte(msg))
	if err != nil {
		logger.LogUM.Error(err)
	}

	err = json.NewEncoder(w).Encode(data)
	if err != nil {
		logger.LogUM.Error(err)
	}
	logger.LogUM.Info(fmt.Sprintf(loggerMessage, code, msg, data))
}

// createErrorResponse create an error response
func createErrorResponse(w http.ResponseWriter, code int, msg string, err error) {
	if err != nil {
		createJSONResponse(w, code, msg, struct {
			Error string `json:"error"`
		}{
			Error: err.Error(),
		})

		return
	}

	createJSONResponse(w, code, msg, nil)
}
