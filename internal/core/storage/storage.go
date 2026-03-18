package storage

import (
	"context"
	"main/internal/features/models"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
)

type DBstruct struct {
	conn *pgx.Conn
}

func Connect(ctx context.Context) (*DBstruct, error) {
	CONN_STRING := os.Getenv("CONN_STRING")
	conn, err := pgx.Connect(ctx, CONN_STRING)
	if err != nil {
		return nil, err
	}
	return &DBstruct{conn: conn}, nil
}

func (s *DBstruct) CreateUser(ctx context.Context, username string, hash_pass []byte) error {
	sqlQuery := ` 
		INSERT INTO users (username,pass_hash,created_at)
		VALUES($1,$2,$3);
	`
	_, err := s.conn.Exec(ctx, sqlQuery, username, hash_pass, time.Now())
	return err
}

func (s *DBstruct) CreateTask(ctx context.Context, userID int, title string, description string, create_time time.Time) (*models.TaskModel, error) {
	sqlQuery := `
		INSERT INTO tasks (user_id, title,description,completed,created_at)
		VALUES($1,$2,$3,$4,$5);
	`
	_, err := s.conn.Exec(ctx, sqlQuery, userID, title, description, false, create_time)

	return &models.TaskModel{Title: title, Description: description, Completed: false, CreatedAt: create_time, CompletedAt: nil}, err
}

func (s *DBstruct) CompleteTask(ctx context.Context, userID, id int) error {
	sqlQuery := `
		UPDATE tasks
		SET completed = $1, completed_at = $2
		WHERE user_id = $3, id = $4
	`
	_, err := s.conn.Exec(ctx, sqlQuery, true, time.Now(), userID, id)
	return err
}

func (s *DBstruct) DeleteTask(ctx context.Context, userID, id int) error {
	sqlQuery := `
		DELETE FROM tasks
		WHERE user_id = $1, id = $2
	`

	_, err := s.conn.Exec(ctx, sqlQuery, userID, id)
	return err
}

func (s *DBstruct) SelectUser(ctx context.Context, username string) (*models.UserModel, error) {
	sqlQuery := `
		SELECT id, username, pass_hash FROM users
		WHERE username = $1
	`
	row := s.conn.QueryRow(ctx, sqlQuery, username)
	var user models.UserModel
	err := row.Scan(
		&user.Id,
		&user.Username,
		&user.PassHash,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *DBstruct) UpdateTask(ctx context.Context, task models.TaskModel) error {
	sqlQuery := `
		UPDATE tasks
		SET title = $1, description = $2, completed = $3, completed_at = $4
		WHERE id = $5
	`
	_, err := s.conn.Exec(ctx, sqlQuery, task.Title, task.Description, task.Completed, task.CompletedAt, task.Id)
	return err
}

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
