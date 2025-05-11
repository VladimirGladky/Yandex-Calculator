package server

import (
	"context"
	task2 "github.com/VladimirGladky/FinalTaskFirstSprint/gen/proto/task"
	"github.com/VladimirGladky/FinalTaskFirstSprint/internal/config"
	"github.com/VladimirGladky/FinalTaskFirstSprint/internal/domain/models"
	"github.com/VladimirGladky/FinalTaskFirstSprint/pkg/calculation"
	"github.com/VladimirGladky/FinalTaskFirstSprint/pkg/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"time"
)

type Agent struct {
	ctx context.Context
	cl  task2.TaskManagementServiceClient
	cfg *config.Config
}

func NewAgent(ctx context.Context, cfg *config.Config) *Agent {
	return &Agent{
		ctx: ctx,
		cfg: cfg,
	}
}

func (a *Agent) Run() {
	conn, err := grpc.NewClient(a.cfg.GrpcHost+":"+a.cfg.GrpcPort, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.GetLoggerFromCtx(a.ctx).Error("error connecting to orchestrator: %v", zap.Error(err))
		return
	}
	defer func(conn *grpc.ClientConn) {
		err = conn.Close()
		if err != nil {
			logger.GetLoggerFromCtx(a.ctx).Error("error closing connection: %v", zap.Error(err))
			return
		}
	}(conn)

	a.cl = task2.NewTaskManagementServiceClient(conn)
	for i := 0; i < a.cfg.ComputingPower; i++ {
		go func() {
			for {
				response, err := a.cl.TaskGet(a.ctx, &task2.TaskGetRequest{})
				if err != nil {
					logger.GetLoggerFromCtx(a.ctx).Error("error getting task: %v", zap.Error(err))
					time.Sleep(15 * time.Second)
					continue
				}
				request := &models.TaskGet{
					Id:            response.Id,
					Arg1:          float64(response.Arg1),
					Arg2:          float64(response.Arg2),
					Operation:     response.Operation,
					OperationTime: int(response.OperationTime),
				}
				res, err4 := calculation.ComputeTask(*request)
				if err4 != nil {
					logger.GetLoggerFromCtx(a.ctx).Error("error computing task: %v", zap.Error(err4))
					time.Sleep(1 * time.Second)
					continue
				}
				logger.GetLoggerFromCtx(a.ctx).Info("result computed", zap.Any("result", res))
				_, err = a.cl.TaskPost(a.ctx, &task2.TaskPostRequest{
					Id:     response.Id,
					Result: float32(res),
				})
				if err != nil {
					logger.GetLoggerFromCtx(a.ctx).Error("error sending result: %v", zap.Error(err))
					return
				}
			}
		}()
	}
	select {}
}
