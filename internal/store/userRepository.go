package store

import "webserver/internal/app/model"

type UserRepository struct {
	store *Store
}

func (r * UserRepository) Create(u *model.User) (*model.User, error) {

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