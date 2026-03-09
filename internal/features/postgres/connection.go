package postgres

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5"
)

func Connect(ctx context.Context) (*DBstruct, error) {
	CONN_STRING := os.Getenv("CONN_STRING")
	conn, err := pgx.Connect(ctx, CONN_STRING)
	if err != nil {
		return nil, err
	}
	return &DBstruct{conn: conn}, nil
}
