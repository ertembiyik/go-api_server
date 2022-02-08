package teststore

import (
	"webserver/internal/app/model"
	"webserver/internal/app/store"
)

type UserRepository struct {
	store *Store
	users map[int]*model.User
}

func (r *UserRepository) Create(u *model.User) error {
	if err := u.Validate(); err != nil {
		return err
	}

	if err := u.BeforeCreate(); err != nil {
		return err
	}

	u.ID = len(r.users) + 1
	r.users[u.ID] = u

	return nil
}

func (r *UserRepository) GetAll() ([]*model.User, error) {
	arr := make([]*model.User, 0, len(r.users))

	for _, user := range r.users {
		arr = append(arr, user)
	}

	return arr, nil
}

func (r *UserRepository) FindByEmail(email string) (*model.User, error) {

	for _, u := range r.users {
		if u.Email == email {
			return u, nil
		}
	}

	return nil, store.ErrorRecordNotFound

}

func (r *UserRepository) Find(id int) (*model.User, error) {
	u, ok := r.users[id]

	if !ok {
		return nil, store.ErrorRecordNotFound
	}

	return u, nil
}

