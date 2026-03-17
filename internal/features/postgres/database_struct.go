package postgres

import (
	"github.com/jackc/pgx/v5"
)

type DBstruct struct {
	conn *pgx.Conn
}
