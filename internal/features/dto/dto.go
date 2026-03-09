package dto

import (
	"encoding/json"
	"errors"
	"time"
)

type UserDTO struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type TaskDTO struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type ErrorDTO struct {
	Message string
	Time    time.Time
}

func (u *UserDTO) ValidateToCreate() error {
	if u.Username == "" {
		return errors.New("username is empty")
	}
	if u.Password == "" {
		return errors.New("password is empty")
	}
	return nil
}

func (t *TaskDTO) ValidateToCreate() error {
	if t.Title == "" {
		return errors.New("Title is empty")
	}
	if t.Description == "" {
		return errors.New("Description is empty")
	}
	return nil
}

func CreateNewError(err error) *ErrorDTO {
	return &ErrorDTO{Message: err.Error(), Time: time.Now()}
}

func (e *ErrorDTO) ToString() string {
	b, err := json.MarshalIndent(e, "", "   ")
	if err != nil {
		panic(err)
	}

	return string(b)
}
