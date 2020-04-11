package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/lvl484/user-manager/logger"
	"github.com/lvl484/user-manager/model"
)

const (
	messageUnauthorized        = "Authenticate failed"
	messageInternalServerError = "Internal server error"
)

func Unauthorized(w http.ResponseWriter) {
	w.Header().Set("WWW-Authenticate", `Basic realm="user-manager"`)
	w.WriteHeader(http.StatusUnauthorized)

	authError := &model.Error{
		Code:    strconv.Itoa(http.StatusUnauthorized),
		Message: messageUnauthorized,
	}

	err := json.NewEncoder(w).Encode(&authError)
	if err != nil {
		logger.LogUM.Errorf("Write Unauthorized response error: %v", err)
	}

	logger.LogUM.Info("Authentication failed! Invalid login or password")
}

func InternalServerError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)

	internalError := &model.Error{
		Code:    strconv.Itoa(http.StatusInternalServerError),
		Message: messageInternalServerError,
	}

	if err := json.NewEncoder(w).Encode(&internalError); err != nil {
		logger.LogUM.Errorf("Write internal server response error: %v", err)
	}

	logger.LogUM.Errorf("Internal server error: %v", err)
}
