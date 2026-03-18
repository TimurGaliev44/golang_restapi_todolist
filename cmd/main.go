package main

import (
	"context"
	"main/internal/core/logger"
	"main/internal/core/storage"
	"main/internal/features/server"
	"main/internal/features/todo"

	"go.uber.org/zap"
)

func main() {
	logLevel := "DEBUG"
	logger, logFileClose, err := logger.NewLogger(logLevel)
	if err != nil {
		panic(err)
	}
	defer logFileClose()

	ctx := context.Background()
	db, err := storage.Connect(ctx)
	if err != nil {
		logger.Error("failed to connect storage", zap.Error(err))
		panic(err)
	}
	rdb, err := storage.NewRedis()
	if err != nil {
		logger.Error("failed to connect redis", zap.Error(err))
	}
	td := todo.New(db, rdb, logger)

	srv := server.NewServer(td)
	if err := srv.StartServer(logger); err != nil {
		panic(err)
	}

}
