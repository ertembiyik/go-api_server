package sqlstore

import (
	"database/sql"
	"webserver/internal/app/model"
	"webserver/internal/app/store"
)

type UserRepository struct {
	store *Store
}

func (r *UserRepository) Create(u *model.User) error {

	if err := u.Validate(); err != nil {
		return err
	}

	if err := u.BeforeCreate(); err != nil {
		return err
	}

	return r.store.db.QueryRow(
		"INSERT INTO users (email, encrypted_password) VALUES ($1, $2) RETURNING id",
		u.Email,
		u.EncryptedPassword).Scan(&u.ID)
	
}

func (r *UserRepository) GetAll() ([]*model.User, error) {
	users := []*model.User{}

	rows, err := r.store.db.Query("SELECT * FROM users ORDER BY id ASC")

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var u *model.User

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
		&u.EncryptedPassword);
		err != nil {
			if err == sql.ErrNoRows {
				return nil, store.ErrorRecordNotFound
			}
		return nil, err
	}

	return u, nil
}
