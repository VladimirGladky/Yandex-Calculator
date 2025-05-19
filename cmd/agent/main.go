package main

import (
	"context"
	"github.com/VladimirGladky/Yandex-Calculator/internal/agent/server"
	"github.com/VladimirGladky/Yandex-Calculator/internal/config"
	"github.com/VladimirGladky/Yandex-Calculator/pkg/logger"
	"go.uber.org/zap"
)

func main() {
	ctx := context.Background()
	ctx, _ = logger.New(ctx)
	cfg, err := config.NewConfig()
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Error("error loading config: %v", zap.Error(err))
		return
	}
	agent := server.NewAgent(ctx, cfg)
	logger.GetLoggerFromCtx(ctx).Info("Agent started")
	agent.Run()
}
