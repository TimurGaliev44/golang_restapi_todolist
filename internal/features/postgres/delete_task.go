package postgres

import (
	"context"
)

func (s *DBstruct) DeleteTask(ctx context.Context, userID, id int) error {
	sqlQuery := `
		DELETE FROM tasks
		WHERE user_id = $1, id = $2
	`

	_, err := s.conn.Exec(ctx, sqlQuery, userID, id)
	return err
}
