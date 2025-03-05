package server

import (
	"encoding/json"
	models2 "github.com/VladimirGladky/FinalTaskFirstSprint/internal/models"
	"github.com/VladimirGladky/FinalTaskFirstSprint/pkg/logger"
	"go.uber.org/zap"
	"io"
	"net/http"
)

func TaskHandlerPost(orchestrator *Orchestrator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				logger.GetLoggerFromCtx(orchestrator.ctx).Error(orchestrator.ctx, "Internal server error")
				w.WriteHeader(http.StatusInternalServerError)
				err0 := json.NewEncoder(w).Encode(models2.BadResponse{Error: "Internal server error1"})
				if err0 != nil {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte(`{"error": "Internal server error2"}`))
				}
				return
			}
		}()
		if r.Method != http.MethodPost {
			//405
			logger.GetLoggerFromCtx(orchestrator.ctx).Error(orchestrator.ctx, "You can use only POST method")
			w.WriteHeader(http.StatusMethodNotAllowed)
			err1 := json.NewEncoder(w).Encode(models2.BadResponse{Error: "You can use only POST method"})
			if err1 != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"error": "Internal server error3"}`))
			}
			return
		}
		request := new(models2.TaskPost)
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				err1 := json.NewEncoder(w).Encode(models2.BadResponse{Error: "Internal server error4"})
				if err1 != nil {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte(`{"error": "Internal server error5"}`))
				}
				return
			}
		}(r.Body)
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			logger.GetLoggerFromCtx(orchestrator.ctx).Error(orchestrator.ctx, "Internal server error")
			w.WriteHeader(http.StatusBadRequest)
			err1 := json.NewEncoder(w).Encode(models2.BadResponse{Error: "Internal server error6"})
			if err1 != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"error": "Internal server error7"}`))
			}
			return
		}
		orchestrator.mu.Lock()
		task, ok := orchestrator.TasksMap[request.Id]
		if !ok {
			logger.GetLoggerFromCtx(orchestrator.ctx).Error(orchestrator.ctx, "Task not found")
			orchestrator.mu.Unlock()
			http.Error(w, `{"error":"Task not found"}`, http.StatusNotFound)
			return
		}
		task.Node.IsLeaf, task.Node.Value = true, request.Result
		delete(orchestrator.TasksMap, request.Id)
		if expression, ex := orchestrator.ExpressionsMap[task.ExprID]; ex {
			orchestrator.SplitTasks(expression)
			if expression.Ast.IsLeaf {
				expression.Status, expression.Result = "completed", expression.Ast.Value
			}
		}
		logger.GetLoggerFromCtx(orchestrator.ctx).Info(orchestrator.ctx, "Task completed", zap.String("task", request.Id), zap.Float64("Result", request.Result))
		orchestrator.mu.Unlock()
		w.WriteHeader(http.StatusOK)
	}
}
