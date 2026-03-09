package postgres

import (
	"context"
	"time"
)

func (s *DBstruct) CreateUser(ctx context.Context, username string, hash_pass []byte) error {
	sqlQuery := ` 
		INSERT INTO users (username,pass_hash,created_at)
		VALUES($1,$2,$3);
	`
	_, err := s.conn.Exec(ctx, sqlQuery, username, hash_pass, time.Now())
	return err
}
