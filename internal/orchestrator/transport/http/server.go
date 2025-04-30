package http

import (
	"context"
	"encoding/json"
	models2 "github.com/VladimirGladky/FinalTaskFirstSprint/internal/models"
	"github.com/VladimirGladky/FinalTaskFirstSprint/internal/orchestrator/service"
	"github.com/VladimirGladky/FinalTaskFirstSprint/pkg/logger"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"io"
	"net/http"
	"os"
)

type Config struct {
	Addr string
}

func ConfigFromEnv() *Config {
	config := new(Config)
	config.Addr = os.Getenv("PORT")
	if config.Addr == "" {
		config.Addr = "4040"
	}
	return config
}

type Orchestrator struct {
	Service *service.Service
	Ctx     context.Context
	config  *Config
}

func New(ctx context.Context, service *service.Service) *Orchestrator {
	return &Orchestrator{
		Ctx:     ctx,
		config:  ConfigFromEnv(),
		Service: service,
	}
}

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (o *Orchestrator) Run() error {
	router := mux.NewRouter()
	router.Use(enableCORS)
	router.HandleFunc("/api/v1/calculate", CalcHandler(o))
	router.HandleFunc("/api/v1/expressions", ExpressionsHandler(o))
	router.HandleFunc("/api/v1/expressions/{id}", ExpressionHandler(o))
	return http.ListenAndServe(":"+o.config.Addr, router)
}

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
		expression, err := orchestrator.Service.CreateExpression(request.Expression)
		if err != nil {
			//400
			logger.GetLoggerFromCtx(orchestrator.Ctx).Error(orchestrator.Ctx, "Bad request", zap.Error(err))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusCreated)
		err4 := json.NewEncoder(w).Encode(models2.ID{ID: expression.Id})
		logger.GetLoggerFromCtx(orchestrator.Ctx).Info(orchestrator.Ctx, "Expression created", zap.String("id", expression.Id))
		if err4 != nil {
			logger.GetLoggerFromCtx(orchestrator.Ctx).Error(orchestrator.Ctx, "Internal server error8")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error": "Internal server error8"}`))
		}
	}
}

func ExpressionHandler(orchestrator *Orchestrator) http.HandlerFunc {
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
		if r.Method != http.MethodGet {
			//405
			logger.GetLoggerFromCtx(orchestrator.Ctx).Error(orchestrator.Ctx, "You can use only GET method")
			w.WriteHeader(http.StatusMethodNotAllowed)
			err10 := json.NewEncoder(w).Encode(models2.BadResponse{Error: "You can use only GET method"})
			if err10 != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"error": "Internal server error3"}`))
			}
			return
		}
		vars := mux.Vars(r)
		id, ok := vars["id"]
		if !ok {
			logger.GetLoggerFromCtx(orchestrator.Ctx).Error(orchestrator.Ctx, "ID is missing")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(models2.BadResponse{Error: "ID is missing"})
			return
		}
		expression, check := orchestrator.Service.GetExpression(id)
		if check != nil {
			logger.GetLoggerFromCtx(orchestrator.Ctx).Error(orchestrator.Ctx, "Expression not found", zap.String("id", id))
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(models2.BadResponse{Error: "Expression not found"})
			return
		}
		err := json.NewEncoder(w).Encode(models2.Expression{Id: expression.Id, Status: expression.Status, Result: expression.Result})
		logger.GetLoggerFromCtx(orchestrator.Ctx).Info(orchestrator.Ctx, "Expression found", zap.String("id", expression.Id), zap.String("status", expression.Status), zap.Float64("result", expression.Result))
		if err != nil {
			logger.GetLoggerFromCtx(orchestrator.Ctx).Error(orchestrator.Ctx, "Internal server error4")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error": "Internal server error4"}`))
			return
		}
	}
}

func ExpressionsHandler(o *Orchestrator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				logger.GetLoggerFromCtx(o.Ctx).Error(o.Ctx, "Internal server error")
				w.WriteHeader(http.StatusInternalServerError)
				err0 := json.NewEncoder(w).Encode(models2.BadResponse{Error: "Internal server error1"})
				if err0 != nil {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte(`{"error": "Internal server error2"}`))
				}
				return
			}
		}()
		if r.Method != http.MethodGet {
			//405
			logger.GetLoggerFromCtx(o.Ctx).Info(o.Ctx, "You can use only GET method")
			w.WriteHeader(http.StatusMethodNotAllowed)
			err10 := json.NewEncoder(w).Encode(models2.BadResponse{Error: "You can use only GET method"})
			if err10 != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"error": "Internal server error3"}`))
			}
			return
		}
		err := json.NewEncoder(w).Encode(models2.Expressions{Expressions: o.Service.GetExpressions()})
		if err != nil {
			logger.GetLoggerFromCtx(o.Ctx).Error(o.Ctx, "Internal server error")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error": "Internal server error4"}`))
			return
		}
		w.WriteHeader(http.StatusOK)
		logger.GetLoggerFromCtx(o.Ctx).Info(o.Ctx, "Expressions sent")
	}
}
