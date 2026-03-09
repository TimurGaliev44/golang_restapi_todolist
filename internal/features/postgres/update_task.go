package postgres

import (
	"context"
	"main/internal/features/models"
)

func (s *DBstruct) UpdateTask(ctx context.Context, task models.TaskModel) error {
	sqlQuery := `
		UPDATE tasks
		SET title = $1, description = $2, completed = $3, completed_at = $4
		WHERE id = $5
	`
	_, err := s.conn.Exec(ctx, sqlQuery, task.Title, task.Description, task.Completed, task.CompletedAt, task.Id)
	return err
}
