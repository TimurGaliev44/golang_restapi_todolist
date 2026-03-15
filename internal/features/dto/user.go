package dto

import "errors"

type UserDTO struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (u *UserDTO) ValidateToCreateUser() error {
	if u.Username == "" {
		return errors.New("username is empty")
	}
	if u.Password == "" {
		return errors.New("password is empty")
	}
	return nil
}
