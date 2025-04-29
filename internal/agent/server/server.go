package server

import (
	"context"
	task2 "github.com/VladimirGladky/FinalTaskFirstSprint/gen/proto/task"
	"github.com/VladimirGladky/FinalTaskFirstSprint/internal/models"
	"github.com/VladimirGladky/FinalTaskFirstSprint/pkg/calculation"
	"github.com/VladimirGladky/FinalTaskFirstSprint/pkg/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"os"
	"strconv"
	"time"
)

type Agent struct {
	ComputingPower int
	ctx            context.Context
	cl             task2.TaskManagementServiceClient
}

func NewAgent(ctx context.Context) *Agent {
	cp, err := strconv.Atoi(os.Getenv("COMPUTING_POWER"))
	if err != nil {
		cp = 1
	}
	return &Agent{
		ComputingPower: cp,
		ctx:            ctx,
	}
}

func (a *Agent) Run() {
	conn, err := grpc.NewClient("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.GetLoggerFromCtx(a.ctx).Error(a.ctx, "error connecting to orchestrator: %v", zap.Error(err))
		return
	}
	defer conn.Close()

	a.cl = task2.NewTaskManagementServiceClient(conn)
	for i := 0; i < a.ComputingPower; i++ {
		go func() {
			for {
				response, err := a.cl.TaskGet(a.ctx, &task2.TaskGetRequest{})
				if err != nil {
					logger.GetLoggerFromCtx(a.ctx).Error(a.ctx, "error getting task: %v", zap.Error(err))
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
					logger.GetLoggerFromCtx(a.ctx).Error(a.ctx, "error computing task: %v", zap.Error(err4))
					time.Sleep(1 * time.Second)
					continue
				}
				logger.GetLoggerFromCtx(a.ctx).Info(a.ctx, "result computed", zap.Any("result", res))
				_, err = a.cl.TaskPost(a.ctx, &task2.TaskPostRequest{
					Id:     response.Id,
					Result: float32(res),
				})
				if err != nil {
					logger.GetLoggerFromCtx(a.ctx).Error(a.ctx, "error sending result: %v", zap.Error(err))
					return
				}
			}
		}()
	}
	select {}
}
