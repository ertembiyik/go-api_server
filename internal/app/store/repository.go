package store

import "webserver/internal/app/model"

type UserRepository interface {
	Create(*model.User) error
	FindByEmail(string) (*model.User, error)
	GetAll() ([]*model.User, error)
	Find(int) (*model.User, error)
}
