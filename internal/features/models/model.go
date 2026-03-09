package models

import "time"

type TaskModel struct {
	Id          int
	User_id     int
	Title       string
	Description string
	Completed   bool
	CreatedAt   time.Time
	CompletedAt *time.Time
}

type UserModel struct {
	Id        int
	Username  string
	PassHash  []byte
	CreatedAt *time.Time
}
