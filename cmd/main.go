package main

import (
	"context"
	"main/internal/core/logger"
	"main/internal/features/postgres"
	"main/internal/features/server"
	"main/internal/features/todo"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}
	if err := os.MkdirAll("out/logs", 0755); err != nil {
		panic(err)
	}
	if err := os.MkdirAll("out/pgdata", 0755); err != nil {
		panic(err)
	}

	logLevel := "DEBUG"
	logger, logFileClose, err := logger.NewLogger(logLevel)
	if err != nil {
		panic(err)
	}
	defer logFileClose()

	ctx := context.Background()
	db, err := postgres.Connect(ctx)
	if err != nil {
		panic(err)
	}
	td := todo.New(db)
	handlers := server.NewHandlers(td, logger)
	srv := server.NewServer(handlers)
	if err := srv.StartServer(); err != nil {
		panic(err)
	}

}
