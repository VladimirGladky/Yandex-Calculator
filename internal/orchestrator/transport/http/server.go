package http

import (
	"context"
	"encoding/json"
	"github.com/VladimirGladky/Yandex-Calculator/internal/config"
	"github.com/VladimirGladky/Yandex-Calculator/internal/domain/models"
	"github.com/VladimirGladky/Yandex-Calculator/internal/orchestrator/service"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"strings"
)

type Orchestrator struct {
	Service service.TaskServiceInterface
	Ctx     context.Context
	Config  *config.Config
}

func New(ctx context.Context, service service.TaskServiceInterface, cfg *config.Config) *Orchestrator {
	return &Orchestrator{
		Ctx:     ctx,
		Config:  cfg,
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
	router.HandleFunc("/api/v1/register", RegisterHandler(o))
	router.HandleFunc("/api/v1/login", LoginHandler(o))
	return http.ListenAndServe(o.Config.OrchestratorHost+":"+o.Config.OrchestratorPort, router)
}

func LoginHandler(orchestrator *Orchestrator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(`{"error": "Internal server error1"}`))
				return
			}
		}()
		if r.Method != http.MethodPost {
			//405
			w.WriteHeader(http.StatusMethodNotAllowed)
			_, _ = w.Write([]byte(`{"error": "You can use only POST method"}`))
			return
		}
		request := new(models.Login)
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(`{"error": "Internal server error4"}`))
			}
		}(r.Body)
		err := json.NewDecoder(r.Body).Decode(request)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`{"error": "Internal server error5"}`))
			return
		}
		token, err := orchestrator.Service.Login(request)
		if err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			err = json.NewEncoder(w).Encode(models.BadResponse{Error: err.Error()})
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(`{"error": "Internal server error6"}`))
			}
			return
		}
		err = json.NewEncoder(w).Encode(models.LoginResponse{Token: token})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`{"error": "Internal server error7"}`))
			return
		}
	}
}

func RegisterHandler(orchestrator *Orchestrator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(`{"error": "Internal server error1"}`))
				return
			}
		}()
		if r.Method != http.MethodPost {
			//405
			w.WriteHeader(http.StatusMethodNotAllowed)
			_, _ = w.Write([]byte(`{"error": "You can use only POST method"}`))
			return
		}
		var req models.RegisterRequest
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(`{"error": "Internal server error4"}`))
				return
			}
		}(r.Body)
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`{"error": "Internal server error5"}`))
			return
		}
		err = orchestrator.Service.Register(&req)
		if err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			err = json.NewEncoder(w).Encode(models.BadResponse{Error: err.Error()})
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(`{"error": "Internal server error6"}`))
			}
			return
		}
		w.WriteHeader(http.StatusOK)
		err14 := json.NewEncoder(w).Encode(models.RegisterGoodResponse{Status: "success"})
		if err14 != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`{"error": "Internal server error10"}`))
			return
		}
	}
}

func CalcHandler(orchestrator *Orchestrator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte(`{"error": "Unauthorized"}`))
			return
		}
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte(`{"error": "Unauthorized"}`))
			return
		}
		tokenString := parts[1]
		defer func() {
			if rec := recover(); rec != nil {
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(`{"error": "Internal server error1"}`))
				return
			}
		}()
		if r.Method != http.MethodPost {
			//405
			w.WriteHeader(http.StatusMethodNotAllowed)
			_, _ = w.Write([]byte(`{"error": "You can use only POST method"}`))
			return
		}
		request := new(models.Request)
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(`{"error": "Internal server error4"}`))
				return
			}
		}(r.Body)
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			//400
			w.WriteHeader(http.StatusBadRequest)
			err2 := json.NewEncoder(w).Encode(models.BadResponse{Error: err.Error()})
			if err2 != nil {
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(`{"error": "Internal server error7"}`))
			}
			return
		}
		expression, err := orchestrator.Service.CreateExpression(request.Expression, tokenString)
		if err != nil {
			//400
			w.WriteHeader(http.StatusBadRequest)
			err3 := json.NewEncoder(w).Encode(models.BadResponse{Error: err.Error()})
			if err3 != nil {
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(`{"error": "Internal server error8"}`))
			}
			return
		}
		w.WriteHeader(http.StatusCreated)
		err4 := json.NewEncoder(w).Encode(models.ID{ID: expression.Id})
		if err4 != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`{"error": "Internal server error9"}`))
		}
	}
}

func ExpressionHandler(orchestrator *Orchestrator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte(`{"error": "Unauthorized"}`))
			return
		}
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte(`{"error": "Unauthorized"}`))
			return
		}
		tokenString := parts[1]
		defer func() {
			if rec := recover(); rec != nil {
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(`{"error": "Internal server error1"}`))
				return
			}
		}()
		if r.Method != http.MethodGet {
			//405
			w.WriteHeader(http.StatusMethodNotAllowed)
			_, _ = w.Write([]byte(`{"error": "You can use only GET method"}`))
			return
		}
		vars := mux.Vars(r)
		id, ok := vars["id"]
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"error": "Id not found"}`))
			return
		}
		expression, check := orchestrator.Service.GetExpression(id, tokenString)
		if check != nil {
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`{"error": "Expression not found"}`))
			return
		}
		err := json.NewEncoder(w).Encode(models.Expression{Id: expression.Id, Status: expression.Status, Result: expression.Result})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`{"error": "Internal server error2"}`))
			return
		}
	}
}

func ExpressionsHandler(o *Orchestrator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte(`{"error": "Unauthorized"}`))
			return
		}
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte(`{"error": "Unauthorized"}`))
			return
		}
		tokenString := parts[1]
		defer func() {
			if rec := recover(); rec != nil {
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(`{"error": "Internal server error1"}`))
				return
			}
		}()
		if r.Method != http.MethodGet {
			//405
			w.WriteHeader(http.StatusMethodNotAllowed)
			_, _ = w.Write([]byte(`{"error": "You can use only GET method"}`))
			return
		}
		expressions, err := o.Service.GetExpressions(tokenString)
		err = json.NewEncoder(w).Encode(models.Expressions{Expressions: expressions})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`{"error": "Internal server error2"}`))
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
