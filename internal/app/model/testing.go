package model

import "testing"

func TestUser(t *testing.T) *User {
	return &User{
		Email: "example@gmail.com",
		Password: "password",
	}
}