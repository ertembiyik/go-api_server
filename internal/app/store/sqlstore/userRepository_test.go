package sqlstore_test

import (
	"testing"
	"webserver/internal/app/model"
	"webserver/internal/app/store"
	"webserver/internal/app/store/sqlstore"
	"github.com/stretchr/testify/assert"
)

func TestUserRepository_Create(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseURL)

	defer teardown("users")

	s := sqlstore.New(db)

	u := model.TestUser(t)

	err := s.User().Create(u)

	assert.NoError(t, err)

	assert.NotNil(t, u)
}

func TestUserRepository_GetAll(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseURL)

	defer teardown("users")

	s := sqlstore.New(db)

	users, err := s.User().GetAll()

	assert.NoError(t, err)

	assert.NotNil(t, users)
}

func TestUserRepository_FindByEmail(t *testing.T) {

	db, teardown := sqlstore.TestDB(t, databaseURL)

	defer teardown("users")

	s := sqlstore.New(db)

	email := "teststor@gmail.com"

	_, err := s.User().FindByEmail(email)

	assert.EqualError(t, err, store.ErrorRecordNotFound.Error())

	u := model.TestUser(t)

	u.Email = email

	s.User().Create(u)

	u, err = s.User().FindByEmail(email)

	assert.NoError(t, err)

	assert.NotNil(t, u)
}
