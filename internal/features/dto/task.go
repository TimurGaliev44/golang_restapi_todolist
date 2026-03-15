package dto

import "errors"

type TaskDTO struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

func (t *TaskDTO) ValidateToCreateTask() error {
	if t.Title == "" {
		return errors.New("title is empty")
	}
	if t.Description == "" {
		return errors.New("description is empty")
	}
	return nil
}
