package apiserver

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"webserver/internal/app/model"
	"webserver/internal/app/store/teststore"

	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/stretchr/testify/assert"
)


func Test_AuthenticateUser(t *testing.T) {

	store := teststore.New()

	u := model.TestUser(t)

	store.User().Create(u)

	testCases := []struct {
		name string
		cookieValue map[interface{}]interface{}
		expectedCode int
	}{
		{
			name: "authenticated",
			cookieValue: map[interface{}]interface{}{
				"user_id": u.ID,
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "not authenticated",
			cookieValue: nil,
			expectedCode: http.StatusUnauthorized,
		},
	}

	secretKey := []byte("secret")
	s := newServer(store, sessions.NewCookieStore(secretKey))
	sc := securecookie.New(secretKey, nil)

	handler := http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(http.StatusOK)
	})

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/", nil)
			cookieStr, _ := sc.Encode(sessionName, tc.cookieValue)
			req.Header.Set("Cookie", fmt.Sprintf("%s=%s", sessionName, cookieStr))
			s.authenticateUser(handler).ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}
func Test_HandleUserCreate(t *testing.T) {

	srv := newServer(teststore.New(), sessions.NewCookieStore([]byte("secret")))

	testCases := []struct {
		name         string
		payload      interface{}
		expectedCode int
	}{
		{
			name: "valid",
			payload: map[string]string{
				"email":    "user@example.org",
				"password": "password",
			},
			expectedCode: http.StatusCreated,
		},
		{
			name:         "invalid payload",
			payload:      "invalid",
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "invalid params",
			payload: map[string]string{
				"email": "invalid",
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			b := &bytes.Buffer{}
			json.NewEncoder(b).Encode(tc.payload)
			req, _ := http.NewRequest(http.MethodPost, "/users", b)
			srv.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}

func Test_HandleSessionCreate(t *testing.T) {

	u := model.TestUser(t)

	store := teststore.New()

	store.User().Create(u)

	srv := newServer(store, sessions.NewCookieStore([]byte("secret")))

	testCases := []struct {
		name         string
		payload      interface{}
		expectedCode int
	}{
		{
			name: "valid",
			payload: map[string]string{
				"email":    u.Email,
				"password": u.Password,
			},
			expectedCode: http.StatusOK,
		},
		{
			name:         "invalid payload",
			payload:      "invalid",
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "invalid email",
			payload: map[string]string{
				"email":    "invalid",
				"password": u.Password,
			},
			expectedCode: http.StatusUnauthorized,
		},
		{
			name: "invalid password",
			payload: map[string]string{
				"email":    u.Email,
				"password": "invalid",
			},
			expectedCode: http.StatusUnauthorized,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			b := &bytes.Buffer{}
			json.NewEncoder(b).Encode(tc.payload)
			req, _ := http.NewRequest(http.MethodPost, "/sessions", b)
			srv.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}
