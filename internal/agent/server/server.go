package server

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/VladimirGladky/FinalTaskFirstSprint/internal/models"
	"github.com/VladimirGladky/FinalTaskFirstSprint/pkg/calculation"
	"github.com/VladimirGladky/FinalTaskFirstSprint/pkg/logger"
	"go.uber.org/zap"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Agent struct {
	ComputingPower int
	ctx            context.Context
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
	for i := 0; i < a.ComputingPower; i++ {
		go func() {
			for {
				req, err1 := http.NewRequest("GET", "http://localhost:4040/internal/task", nil)
				if err1 != nil {
					logger.GetLoggerFromCtx(a.ctx).Error(a.ctx, "error getting task: %v", zap.Error(err1))
					time.Sleep(1 * time.Second)
					continue
				}
				resp, err2 := http.DefaultClient.Do(req)
				if err2 != nil {
					logger.GetLoggerFromCtx(a.ctx).Error(a.ctx, "error getting task: %v", zap.Error(err2))
					time.Sleep(1 * time.Second)
					continue
				}
				if resp.StatusCode == http.StatusNotFound {
					logger.GetLoggerFromCtx(a.ctx).Info(a.ctx, "no task")
					time.Sleep(15 * time.Second)
					continue
				}
				request := new(models.TaskGet)
				err3 := json.NewDecoder(resp.Body).Decode(&request)
				resp.Body.Close()
				if err3 != nil {
					logger.GetLoggerFromCtx(a.ctx).Error(a.ctx, "error decoding task: %v", zap.Error(err3))
					time.Sleep(1 * time.Second)
					continue
				}
				res, err4 := calculation.ComputeTask(*request)
				if err4 != nil {
					logger.GetLoggerFromCtx(a.ctx).Error(a.ctx, "error computing task: %v", zap.Error(err4))
					time.Sleep(1 * time.Second)
					continue
				}
				jsonData, err5 := json.Marshal(models.TaskPost{Id: request.Id, Result: res})
				if err5 != nil {
					logger.GetLoggerFromCtx(a.ctx).Error(a.ctx, "error marshaling result: %v", zap.Error(err5))
					return
				}
				logger.GetLoggerFromCtx(a.ctx).Info(a.ctx, "result sent", zap.String("id", request.Id), zap.Float64("result", res))
				req, err := http.NewRequest("POST", "http://localhost:4040/internal/task", bytes.NewBuffer(jsonData))
				if err != nil {
					logger.GetLoggerFromCtx(a.ctx).Error(a.ctx, "error sending result: %v", zap.Error(err))
					continue
				}
				resp, err = http.DefaultClient.Do(req)
				if err != nil {
					logger.GetLoggerFromCtx(a.ctx).Error(a.ctx, "error sending result: %v", zap.Error(err))
					continue
				}
				resp.Body.Close()
			}
		}()
	}
	select {}
}
