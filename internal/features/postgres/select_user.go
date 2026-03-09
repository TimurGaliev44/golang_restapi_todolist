package postgres

import (
	"context"
	"main/internal/features/models"
)

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
