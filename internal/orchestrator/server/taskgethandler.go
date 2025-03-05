package server

import (
	"encoding/json"
	models2 "github.com/VladimirGladky/FinalTaskFirstSprint/internal/models"
	"github.com/VladimirGladky/FinalTaskFirstSprint/pkg/logger"
	"go.uber.org/zap"
	"net/http"
)

func TaskHandlerGet(orchestrator *Orchestrator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				logger.GetLoggerFromCtx(orchestrator.ctx).Error(orchestrator.ctx, "Internal server error")
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(models2.BadResponse{Error: "Internal server error1"})
				return
			}
		}()
		if r.Method != http.MethodGet {
			//405
			logger.GetLoggerFromCtx(orchestrator.ctx).Info(orchestrator.ctx, "You can use only GET method")
			w.WriteHeader(http.StatusMethodNotAllowed)
			json.NewEncoder(w).Encode(models2.BadResponse{Error: "You can use only GET method"})
			return
		}
		if len(orchestrator.TasksArr) == 0 {
			logger.GetLoggerFromCtx(orchestrator.ctx).Info(orchestrator.ctx, "No task available")
			http.Error(w, `{"error":"No task available"}`, http.StatusNotFound)
			return
		}
		task := orchestrator.TasksArr[0]
		orchestrator.TasksArr = orchestrator.TasksArr[1:]
		if expr, exists := orchestrator.ExpressionsMap[task.ExprID]; exists {
			expr.Status = "in_progress"
		}
		err := json.NewEncoder(w).Encode(task)
		logger.GetLoggerFromCtx(orchestrator.ctx).Info(orchestrator.ctx, "Task sent", zap.String("task", task.ID), zap.Float64("Arg1", task.Arg1), zap.Float64("Arg2", task.Arg2), zap.String("Operation", task.Operation))
		if err != nil {
			logger.GetLoggerFromCtx(orchestrator.ctx).Error(orchestrator.ctx, "Internal server error")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(models2.BadResponse{Error: "Internal server error2"})
			return
		}
	}
}
