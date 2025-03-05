package main

import (
	"context"
	"github.com/VladimirGladky/FinalTaskFirstSprint/internal/agent/server"
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
	agent := server.NewAgent(ctx)
	logger.GetLoggerFromCtx(ctx).Info(ctx, "Agent started")
	agent.Run()
}
