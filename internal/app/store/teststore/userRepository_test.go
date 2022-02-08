package teststore_test

import (
	"testing"
	"webserver/internal/app/model"
	"webserver/internal/app/store"
	"webserver/internal/app/store/teststore"

	"github.com/stretchr/testify/assert"
)

func TestUserRepository_Create(t *testing.T) {

	s := teststore.New()

	u := model.TestUser(t)

	err := s.User().Create(u)

	assert.NoError(t, err)

	assert.NotNil(t, u)
}

func TestUserRepository_GetAll(t *testing.T) {
	s := teststore.New()

	users, err := s.User().GetAll()

	assert.NoError(t, err)

	assert.NotNil(t, users)
}

func TestUserRepository_FindByEmail(t *testing.T) {
	s := teststore.New()

	email := "teststore@gmail.com"

	_, err := s.User().FindByEmail(email)

	assert.EqualError(t, err, store.ErrorRecordNotFound.Error())

	u := model.TestUser(t)

	u.Email = email

	s.User().Create(u)

	u, err = s.User().FindByEmail(email)

	assert.NoError(t, err)

	assert.NotNil(t, u)
}


func TestUserRepository_Find(t *testing.T) {

	s := teststore.New()

	u := model.TestUser(t)

	_, err := s.User().Find(u.ID)

	assert.EqualError(t, err, store.ErrorRecordNotFound.Error())

	s.User().Create(u)

	u, err = s.User().Find(u.ID)

	assert.NoError(t, err)

	assert.NotNil(t, u)
}
