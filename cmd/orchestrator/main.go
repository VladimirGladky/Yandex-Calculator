package main

import (
	"context"
	"fmt"
	"github.com/VladimirGladky/FinalTaskFirstSprint/gen/proto/task"
	service2 "github.com/VladimirGladky/FinalTaskFirstSprint/internal/orchestrator/service"
	gr "github.com/VladimirGladky/FinalTaskFirstSprint/internal/orchestrator/transport/grpc"
	"github.com/VladimirGladky/FinalTaskFirstSprint/internal/orchestrator/transport/http"
	"github.com/VladimirGladky/FinalTaskFirstSprint/pkg/logger"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	wg := &sync.WaitGroup{}
	ctx, _ = logger.New(ctx)
	_ = godotenv.Load("local.env")
	lis, err := net.Listen("tcp", fmt.Sprintf("%s", "localhost:8080"))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	srv := service2.NewService()
	orch := http.New(ctx, srv)
	service := gr.NewService(orch)
	grpcServer := grpc.NewServer()
	task.RegisterTaskManagementServiceServer(grpcServer, service)
	wg.Add(1)
	go func() {
		defer wg.Done()
		logger.GetLoggerFromCtx(ctx).Info(ctx, "Orchestrator started")
		if err := orch.Run(); err != nil {
			logger.GetLoggerFromCtx(ctx).Error(ctx, "Orchestrator error", zap.Error(err))
			cancel()
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := grpcServer.Serve(lis); err != nil {
			logger.GetLoggerFromCtx(ctx).Error(ctx, "Orchestrator error", zap.Error(err))
			cancel()
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-sigCh:
		logger.GetLoggerFromCtx(ctx).Info(ctx, "Forced shutdown")
		grpcServer.GracefulStop()
		cancel()
	case <-ctx.Done():
		logger.GetLoggerFromCtx(ctx).Info(ctx, "Graceful shutdown")
	}

	wg.Wait()
}
