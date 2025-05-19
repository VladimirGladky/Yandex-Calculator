package main

import (
	"context"
	"github.com/VladimirGladky/Yandex-Calculator/internal/app"
	"github.com/VladimirGladky/Yandex-Calculator/internal/config"
	"github.com/VladimirGladky/Yandex-Calculator/pkg/logger"
	"go.uber.org/zap"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	cfg, err := config.NewConfig()
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Error("error loading config: %v", zap.Error(err))
		return
	}
	ctx, _ = logger.New(ctx)
	application := app.New(cfg, ctx)
	application.MustRun()
}
