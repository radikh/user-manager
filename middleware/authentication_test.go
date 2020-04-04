package middleware

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/lvl484/user-manager/logger"
	"github.com/lvl484/user-manager/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	logger.SetLogger(&logger.LogConfig{Output: "Stdout", Level: "debug"})

	code := m.Run()

	os.Exit(code)
}

func TestBasicAuthenticationMiddlewareValidPass(t *testing.T) {
	ma := &MockAuth{}
	ba := NewBasicAuthentication(ma)

	r, err := http.NewRequest("GET", "/summer", nil)
	require.NoError(t, err)

	r.SetBasicAuth("i3odja", "1q2w3e4r")

	w := httptest.NewRecorder()

	ba.Middleware(OK).ServeHTTP(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestBasicAuthenticationMiddlewareInvalidPass(t *testing.T) {
	ma := &MockAuth{}
	ba := NewBasicAuthentication(ma)

	r, err := http.NewRequest("GET", "/summer", nil)
	require.NoError(t, err)

	r.SetBasicAuth("i3odja", "123123")

	w := httptest.NewRecorder()

	ba.Middleware(OK).ServeHTTP(w, r)

	checkErrorResponse(t, w, http.StatusUnauthorized)
}

func TestBasicAuthenticationMiddlewareError(t *testing.T) {
	ma := &MockAuth{Error: errors.New("middleware error")}
	ba := NewBasicAuthentication(ma)

	r, err := http.NewRequest("GET", "/summer", nil)
	require.NoError(t, err)

	r.SetBasicAuth("i3odja", "123123")

	w := httptest.NewRecorder()

	ba.Middleware(OK).ServeHTTP(w, r)

	checkErrorResponse(t, w, http.StatusInternalServerError)
}

func checkErrorResponse(t *testing.T, w *httptest.ResponseRecorder, expectedCode int) {
	mError := model.Error{}

	err := json.Unmarshal(w.Body.Bytes(), &mError)
	require.NoError(t, err)

	assert.NotEmpty(t, mError.Code)
	assert.NotEmpty(t, mError.Message)

	assert.Equal(t, expectedCode, w.Code)
}

type MockAuth struct {
	Error error
}

func (ma *MockAuth) GetInfo(username string) (*model.User, error) {
	return &model.User{
		ID:        "123e4567-e89b-12d3-a456-426655440000",
		Username:  username,
		Password:  "$argon2id$v=19$m=65536,t=3,p=1$Ga9X4EvymyOzUoz+uVMy6w$y5sQVQ",
		Email:     "qwerty@gmail.com",
		FirstName: "UserF",
		LastName:  "UserL",
		Phone:     0671112233,
		CreatedAt: nil,
		UpdatedAt: nil,
	}, ma.Error
}

var OK http.Handler = http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusOK)
})
