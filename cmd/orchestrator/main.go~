package main

import (
	"context"
	"github.com/VladimirGladky/FinalTaskFirstSprint/internal/orchestrator/server"
	"github.com/VladimirGladky/FinalTaskFirstSprint/pkg/logger"
	"github.com/joho/godotenv"
)

func main() {
	ctx := context.Background()
	ctx, _ = logger.New(ctx)
	err := godotenv.Load("local.env")
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Error(ctx, "Error loading .env file")
	}
	srv := server.New(ctx)
	logger.GetLoggerFromCtx(ctx).Info(ctx, "Orchestrator started")
	srv.Run()
}
