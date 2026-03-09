package postgres

import (
	"context"
	"main/internal/features/models"
)

func (s *DBstruct) SelectAll(ctx context.Context, userID int) ([]models.TaskModel, error) {
	sqlQuery := `
		SELECT id, title, description, completed, created_at, completed_at
		FROM tasks
		WHERE user_id = $1
	`

	rows, err := s.conn.Query(ctx, sqlQuery, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := make([]models.TaskModel, 0)
	for rows.Next() {
		var task models.TaskModel
		err = rows.Scan(
			&task.Id,
			&task.Title,
			&task.Description,
			&task.Completed,
			&task.CreatedAt,
			&task.CompletedAt,
		)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (s *DBstruct) SelectID(ctx context.Context, userID int, id int) (*models.TaskModel, error) {
	sqlQuery := `
		SELECT id, title, description, completed, created_at, completed_at
		FROM tasks
		WHERE user_id = $1
	`
	row := s.conn.QueryRow(ctx, sqlQuery, userID)
	var task models.TaskModel
	err := row.Scan(
		&task.Id,
		&task.User_id,
		&task.Title,
		&task.Description,
		&task.Completed,
		&task.CreatedAt,
		&task.CompletedAt,
	)
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (s *DBstruct) SelectUncompleted(ctx context.Context, userID int) ([]models.TaskModel, error) {
	sqlQuery := `
		SELECT id, title, description, completed, created_at, completed_at
		FROM tasks
		WHERE user_id = $1, completed = false
	`

	rows, err := s.conn.Query(ctx, sqlQuery, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := make([]models.TaskModel, 0)
	for rows.Next() {
		var task models.TaskModel
		err = rows.Scan(
			&task.Id,
			&task.User_id,
			&task.Title,
			&task.Description,
			&task.Completed,
			&task.CreatedAt,
			&task.CompletedAt,
		)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}
