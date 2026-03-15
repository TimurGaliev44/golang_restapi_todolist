package dto

import "errors"

type TaskDTO struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

func (u *UserDTO) ValidateToCreateTask() error {
	if u.Username == "" {
		return errors.New("username is empty")
	}
	if u.Password == "" {
		return errors.New("password is empty")
	}
	return nil
}
