package grpcapp

import (
	"context"
	"fmt"
	"github.com/VladimirGladky/FinalTaskFirstSprint/internal/config"
	"github.com/VladimirGladky/FinalTaskFirstSprint/internal/orchestrator/service"
	"github.com/VladimirGladky/FinalTaskFirstSprint/internal/orchestrator/transport/taskgRPC"
	"github.com/VladimirGladky/FinalTaskFirstSprint/pkg/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
)

type App struct {
	GRPCSrv *grpc.Server
	cfg     *config.Config
	ctx     context.Context
}

func New(cfg *config.Config, service *service.Service, ctx context.Context) *App {
	gRPC := grpc.NewServer()
	taskgRPC.Register(gRPC, service)

	return &App{
		GRPCSrv: gRPC,
		cfg:     cfg,
		ctx:     ctx,
	}
}

func (a *App) Run() error {
	lis, err := net.Listen("tcp", fmt.Sprintf("%s", a.cfg.GrpcHost+":"+a.cfg.GrpcPort))
	if err != nil {
		logger.GetLoggerFromCtx(a.ctx).Error("error listening: %v", zap.Error(err))
		return err
	}
	if err = a.GRPCSrv.Serve(lis); err != nil {
		logger.GetLoggerFromCtx(a.ctx).Error("error serving: %v", zap.Error(err))
		return err
	}
	return nil
}
