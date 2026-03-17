package main

import (
	"context"
	"main/internal/core/logger"
	"main/internal/features/postgres"
	"main/internal/features/server"
	"main/internal/features/todo"
	"os"

	"go.uber.org/zap"
)

func main() {
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
		logger.Error("failed to connect storage", zap.Error(err))
		panic(err)
	}
	td := todo.New(db)

	srv := server.NewServer(td, logger)
	if err := srv.StartServer(); err != nil {
		panic(err)
	}

}
