package server

import (
	"encoding/json"
	"github.com/VladimirGladky/FinalTaskFirstSprint/internal/models"
	"github.com/VladimirGladky/FinalTaskFirstSprint/pkg/logger"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"net/http"
)

func ExpressionHandler(orchestrator *Orchestrator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				logger.GetLoggerFromCtx(orchestrator.ctx).Error(orchestrator.ctx, "Internal server error1")
				w.WriteHeader(http.StatusInternalServerError)
				err0 := json.NewEncoder(w).Encode(models.BadResponse{Error: "Internal server error1"})
				if err0 != nil {
					logger.GetLoggerFromCtx(orchestrator.ctx).Error(orchestrator.ctx, "Internal server error2")
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte(`{"error": "Internal server error2"}`))
				}
				return
			}
		}()
		if r.Method != http.MethodGet {
			//405
			logger.GetLoggerFromCtx(orchestrator.ctx).Error(orchestrator.ctx, "You can use only GET method")
			w.WriteHeader(http.StatusMethodNotAllowed)
			err10 := json.NewEncoder(w).Encode(models.BadResponse{Error: "You can use only GET method"})
			if err10 != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"error": "Internal server error3"}`))
			}
			return
		}
		vars := mux.Vars(r)
		id, ok := vars["id"]
		if !ok {
			logger.GetLoggerFromCtx(orchestrator.ctx).Error(orchestrator.ctx, "ID is missing")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(models.BadResponse{Error: "ID is missing"})
			return
		}
		expression, ok := orchestrator.ExpressionsMap[id]
		if !ok {
			logger.GetLoggerFromCtx(orchestrator.ctx).Error(orchestrator.ctx, "Expression not found")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(models.BadResponse{Error: "Expression not found"})
			return
		}
		err := json.NewEncoder(w).Encode(models.Expression{Id: expression.Id, Status: expression.Status, Result: expression.Result})
		logger.GetLoggerFromCtx(orchestrator.ctx).Info(orchestrator.ctx, "Expression found", zap.String("id", expression.Id), zap.String("status", expression.Status), zap.Float64("result", expression.Result))
		if err != nil {
			logger.GetLoggerFromCtx(orchestrator.ctx).Error(orchestrator.ctx, "Internal server error4")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error": "Internal server error4"}`))
			return
		}
	}
}
