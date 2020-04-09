package middleware_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/lvl484/user-manager/middleware"
	"github.com/lvl484/user-manager/mock"

	gomock "github.com/golang/mock/gomock"
	"github.com/lvl484/user-manager/logger"
	model "github.com/lvl484/user-manager/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	logger.SetLogger(&logger.LogConfig{Output: "Stdout", Level: "debug"})

	code := m.Run()

	os.Exit(code)
}

func TestBasicAuthenticationMiddlewareValidPass(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock := mock.NewMockUserProvider(ctrl)

	mock.EXPECT().GetInfo("i3odja").Return(userInfo, nil)

	ba := middleware.NewBasicAuthentication(mock)

	r, err := http.NewRequest("GET", "/summer", nil)
	require.NoError(t, err)

	r.SetBasicAuth("i3odja", "1q2w3e4r")

	w := httptest.NewRecorder()

	ba.Middleware(wrappedHandler).ServeHTTP(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestBasicAuthenticationMiddlewareInvalidPass(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock := mock.NewMockUserProvider(ctrl)

	mock.EXPECT().GetInfo("i3odja").Return(userInfo, nil)

	ba := middleware.NewBasicAuthentication(mock)

	r, err := http.NewRequest("GET", "/summer", nil)
	require.NoError(t, err)

	r.SetBasicAuth("i3odja", "123123")

	w := httptest.NewRecorder()

	ba.Middleware(wrappedHandler).ServeHTTP(w, r)

	checkErrorResponse(t, w, http.StatusUnauthorized)
}

func TestBasicAuthenticationMiddlewareError(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock := mock.NewMockUserProvider(ctrl)

	mock.EXPECT().GetInfo("i3odja").Return(nil, errors.New("middleware error"))

	ba := middleware.NewBasicAuthentication(mock)

	r, err := http.NewRequest("GET", "/summer", nil)
	require.NoError(t, err)

	r.SetBasicAuth("i3odja", "123123")

	w := httptest.NewRecorder()

	ba.Middleware(wrappedHandler).ServeHTTP(w, r)

	checkErrorResponse(t, w, http.StatusInternalServerError)
}

func TestBasicAuthenticationMiddlewareResponseInvalid(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock := mock.NewMockUserProvider(ctrl)

	ba := middleware.NewBasicAuthentication(mock)

	r, err := http.NewRequest("GET", "/summer", nil)
	require.NoError(t, err)

	w := httptest.NewRecorder()

	ba.Middleware(wrappedHandler).ServeHTTP(w, r)

	checkErrorResponse(t, w, http.StatusUnauthorized)
}

var wrappedHandler http.Handler = http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusOK)
})

func checkErrorResponse(t *testing.T, w *httptest.ResponseRecorder, expectedCode int) {
	mError := model.Error{}

	err := json.Unmarshal(w.Body.Bytes(), &mError)
	require.NoError(t, err)

	filename := fmt.Sprintf("./testdata/response_%d.json", expectedCode)
	expected, err := ioutil.ReadFile(filename)
	require.NoError(t, err)

	assert.JSONEq(t, string(expected), w.Body.String())

	assert.Equal(t, expectedCode, w.Code)
}

var userInfo = &model.User{
	ID:        "123e4567-e89b-12d3-a456-426655440000",
	Username:  "i3odja",
	Password:  "$argon2id$v=19$m=65536,t=3,p=1$Ga9X4EvymyOzUoz+uVMy6w$y5sQVQ",
	Email:     "qwerty@gmail.com",
	FirstName: "UserF",
	LastName:  "UserL",
	Phone:     "0671112233",
}
