package postgres

import (
	"context"
	"main/internal/features/models"
	"time"
)

func (s *DBstruct) CreateTask(ctx context.Context, userID int, title string, description string, create_time time.Time) (*models.TaskModel, error) {
	sqlQuery := `
		INSERT INTO tasks (user_id, title,description,completed,created_at)
		VALUES($1,$2,$3,$4,$5);
	`
	_, err := s.conn.Exec(ctx, sqlQuery, userID, title, description, false, create_time)

	return &models.TaskModel{Title: title, Description: description, Completed: false, CreatedAt: create_time, CompletedAt: nil}, err
}
