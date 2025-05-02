package app

import (
	"context"
	grpcapp "github.com/VladimirGladky/FinalTaskFirstSprint/internal/app/grpc"
	orchestratorapp "github.com/VladimirGladky/FinalTaskFirstSprint/internal/app/orchestrator"
	"github.com/VladimirGladky/FinalTaskFirstSprint/internal/config"
	"github.com/VladimirGladky/FinalTaskFirstSprint/internal/orchestrator/service"
	"github.com/VladimirGladky/FinalTaskFirstSprint/internal/orchestrator/transport/http"
	"github.com/VladimirGladky/FinalTaskFirstSprint/pkg/logger"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type App struct {
	GRPCsrv      *grpcapp.App
	Orchestrator *orchestratorapp.App
	ctx          context.Context
	wg           sync.WaitGroup
	cancel       context.CancelFunc
}

func New(
	cfg *config.Config,
	ctx context.Context,
	srv *service.Service,
) *App {
	orchestrator := http.New(ctx, srv, cfg)
	orchApp := orchestratorapp.New(orchestrator)
	grpcApp := grpcapp.New(cfg, orchestrator, ctx)
	return &App{
		GRPCsrv:      grpcApp,
		Orchestrator: orchApp,
		ctx:          ctx,
	}
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Run() error {
	errCh := make(chan error, 2)
	a.wg.Add(1)
	go func() {
		logger.GetLoggerFromCtx(a.ctx).Info("Orchestrator started")
		defer a.wg.Done()
		if err := a.Orchestrator.Run(); err != nil {
			errCh <- err
			a.cancel()
		}
	}()
	a.wg.Add(1)
	go func() {
		defer a.wg.Done()
		logger.GetLoggerFromCtx(a.ctx).Info("gRPC server started")
		if err := a.GRPCsrv.Run(); err != nil {
			errCh <- err
			a.cancel()
		}
	}()
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	select {
	case err := <-errCh:
		logger.GetLoggerFromCtx(a.ctx).Error("error running app", zap.Error(err))
		return err
	case sig := <-sigCh:
		logger.GetLoggerFromCtx(a.ctx).Info("received signal", zap.String("signal", sig.String()))
		a.Stop()
	case <-a.ctx.Done():
		logger.GetLoggerFromCtx(a.ctx).Info("context done")
	}

	return nil
}

func (a *App) Stop() {
	logger.GetLoggerFromCtx(a.ctx).Info("stopping app")
	a.cancel()
	a.GRPCsrv.GRPCsrv.GracefulStop()
	a.wg.Wait()
	logger.GetLoggerFromCtx(a.ctx).Info("app stopped")
}
