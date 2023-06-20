package models

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       int    `json:"id"`
	Username string `json:"user_name"`
	Password string `json:"password"`
	IsAdmin  bool   `json:"is_admin"`
}

func (u *User) PasswordMatches(plainText string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(plainText))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword): // invalid password
			return false, nil
		default:
			return false, err
		}
	}
	return true, nil
}
