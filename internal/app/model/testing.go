package model

import "testing"

func TestUser(t *testing.T) *User {
	return &User{
		Email: "e@gmail.com",
		Password: "password",
	}
}