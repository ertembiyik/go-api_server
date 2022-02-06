package store_test

import (
	"testing"
	"webserver/internal/app/model"
	"webserver/internal/store"
	"github.com/stretchr/testify/assert"
)

func TestUserRepository_Create(t *testing.T) {
	s, teardown := store.TestStore(t, databaseURL) 

	defer teardown("users")

	u, err := s.User().Create(model.TestUser(t))

	assert.NoError(t, err)

	assert.NotNil(t, u)
}	

func TestUserRepository_GetAll(t *testing.T) {
	s, teardown := store.TestStore(t, databaseURL) 

	defer teardown("users")

	users, err := s.User().GetAll()

	assert.NoError(t, err)
	assert.NotNil(t, users)
}