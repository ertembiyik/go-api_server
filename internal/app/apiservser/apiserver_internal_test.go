package apiserver

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)


func TestAPIServer_HandleGetNotes(t *testing.T) {
	s := New(NewConfig())
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/users", nil)

	s.getUsers().ServeHTTP(rec, req)

	assert.Equal(t, rec.Body.String(), "Hello")
}