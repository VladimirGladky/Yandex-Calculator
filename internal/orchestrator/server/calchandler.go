package server

import (
	"encoding/json"
	models2 "github.com/VladimirGladky/FinalTaskFirstSprint/internal/models"
	"github.com/VladimirGladky/FinalTaskFirstSprint/internal/orchestrator/parser"
	"github.com/VladimirGladky/FinalTaskFirstSprint/pkg/logger"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"io"
	"net/http"
)

func CalcHandler(orchestrator *Orchestrator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				logger.GetLoggerFromCtx(orchestrator.Ctx).Error(orchestrator.Ctx, "Internal server error1")
				w.WriteHeader(http.StatusInternalServerError)
				err0 := json.NewEncoder(w).Encode(models2.BadResponse{Error: "Internal server error1"})
				if err0 != nil {
					logger.GetLoggerFromCtx(orchestrator.Ctx).Error(orchestrator.Ctx, "Internal server error2")
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte(`{"error": "Internal server error2"}`))
				}
				return
			}
		}()
		if r.Method != http.MethodPost {
			//405
			logger.GetLoggerFromCtx(orchestrator.Ctx).Error(orchestrator.Ctx, "You can use only POST method")
			w.WriteHeader(http.StatusMethodNotAllowed)
			err10 := json.NewEncoder(w).Encode(models2.BadResponse{Error: "You can use only POST method"})
			if err10 != nil {
				logger.GetLoggerFromCtx(orchestrator.Ctx).Error(orchestrator.Ctx, "Internal server error3")
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"error": "Internal server error3"}`))
			}
			return
		}
		request := new(models2.Request)
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
			//400
			logger.GetLoggerFromCtx(orchestrator.Ctx).Error(orchestrator.Ctx, "Bad request", zap.Error(err))
			w.WriteHeader(http.StatusBadRequest)
			err2 := json.NewEncoder(w).Encode(models2.BadResponse{Error: "Bad request"})
			if err2 != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"error": "Internal server error6"}`))
			}
			return
		}
		id := uuid.New().String()
		ast, err := parser.BuildExpressionTree(request.Expression)
		if err != nil {
			//422
			logger.GetLoggerFromCtx(orchestrator.Ctx).Error(orchestrator.Ctx, "Expression is not valid", zap.Error(err))
			w.WriteHeader(http.StatusUnprocessableEntity)
			err3 := json.NewEncoder(w).Encode(models2.BadResponse{Error: "Expression is not valid"})
			if err3 != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"error": "Internal server error7"}`))
			}
			return
		}
		orchestrator.Mu.Lock()
		orchestrator.ExpressionsMap[id] = &models2.Expression{Id: id, Status: "in progress", Ast: ast}
		orchestrator.SplitTasks(orchestrator.ExpressionsMap[id])
		orchestrator.Mu.Unlock()

		w.WriteHeader(http.StatusCreated)
		err4 := json.NewEncoder(w).Encode(models2.ID{ID: id})
		logger.GetLoggerFromCtx(orchestrator.Ctx).Info(orchestrator.Ctx, "Expression created", zap.String("id", id))
		if err4 != nil {
			logger.GetLoggerFromCtx(orchestrator.Ctx).Error(orchestrator.Ctx, "Internal server error8")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error": "Internal server error8"}`))
		}
	}
}
