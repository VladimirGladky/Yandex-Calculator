package grpcapp

import (
	"context"
	"fmt"
	"github.com/VladimirGladky/FinalTaskFirstSprint/internal/config"
	"github.com/VladimirGladky/FinalTaskFirstSprint/internal/orchestrator/transport/http"
	"github.com/VladimirGladky/FinalTaskFirstSprint/internal/orchestrator/transport/taskgRPC"
	"github.com/VladimirGladky/FinalTaskFirstSprint/pkg/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
)

type App struct {
	GRPCsrv *grpc.Server
	cfg     *config.Config
	ctx     context.Context
	orch    *http.Orchestrator
}

func New(cfg *config.Config, orchestrator *http.Orchestrator, ctx context.Context) *App {
	gRPC := grpc.NewServer()
	taskgRPC.Register(gRPC, orchestrator)

	return &App{
		GRPCsrv: gRPC,
		cfg:     cfg,
		ctx:     ctx,
		orch:    orchestrator,
	}
}

func (a *App) Run() error {
	lis, err := net.Listen("tcp", fmt.Sprintf("%s", a.cfg.GrpcHost+":"+a.cfg.GrpcPort))
	if err != nil {
		logger.GetLoggerFromCtx(a.ctx).Error("error listening: %v", zap.Error(err))
		return err
	}
	if err = a.GRPCsrv.Serve(lis); err != nil {
		logger.GetLoggerFromCtx(a.ctx).Error("error serving: %v", zap.Error(err))
		return err
	}
	return nil
}
