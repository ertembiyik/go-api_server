package store

import (
	"webserver/internal/app/model"
)

type UserRepository struct {
	store *Store
}

func (r *UserRepository) Create(u *model.User) (*model.User, error) {

	if err := u.Validate(); err != nil {
		return nil, err
	}

	if err := u.BeforeCreate(); err != nil {
		return nil, err
	}

	if err := r.store.db.QueryRow(
		"INSERT INTO users (email, encrypted_password) VALUES ($1, $2) RETURNING id",
		u.Email,
		u.Password).Scan(&u.ID); err != nil {
		return nil, err
	}

	return u, nil
}

func (r *UserRepository) GetAll() ([]model.User, error) {
	users := []model.User{}

	rows, err := r.store.db.Query("SELECT * FROM users ORDER BY id ASC")

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var u model.User

		if err := rows.Scan(&u.ID, &u.Email, &u.Password); err != nil {
			return nil, err
		}

		users = append(users, u)
	}

	return users, nil
}

func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
	u := &model.User{}

	if err := r.store.db.QueryRow("SELECT id, email, encrypted_password FROM users WHERE email = $1", email).Scan(
		&u.ID,
		&u.Email,
		&u.EncryptedPassword); err != nil {
		return nil, err
	}

	return u, nil
}
