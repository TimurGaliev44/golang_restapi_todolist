package postgres

import (
	"context"
	"time"
)

func (s *DBstruct) CompleteTask(ctx context.Context, userID, id int) error {
	sqlQuery := `
		UPDATE tasks
		SET completed = $1, completed_at = $2
		WHERE user_id = $3, id = $4
	`
	_, err := s.conn.Exec(ctx, sqlQuery, true, time.Now(), userID, id)
	return err
}
