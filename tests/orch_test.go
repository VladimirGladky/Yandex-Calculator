package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/VladimirGladky/FinalTaskFirstSprint/internal/config"
	"github.com/VladimirGladky/FinalTaskFirstSprint/internal/domain/models"
	http2 "github.com/VladimirGladky/FinalTaskFirstSprint/internal/orchestrator/transport/http"
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockTaskService struct {
	registerFunc       func(*models.RegisterRequest) error
	loginFunc          func(*models.Login) (string, error)
	createExprFunc     func(string, string) (*models.Expression, error)
	getExpressionFunc  func(string, string) (*models.Expression, error)
	getExpressionsFunc func(string) ([]*models.Expression, error)
	splitTasksFunc     func(*models.Expression)
}

func (m *mockTaskService) GetExpression(id string, token string) (*models.Expression, error) {
	return m.getExpressionFunc(id, token)
}

func (m *mockTaskService) CreateExpression(expr string, token string) (*models.Expression, error) {
	return m.createExprFunc(expr, token)
}

func (m *mockTaskService) GetExpressions(token string) ([]*models.Expression, error) {
	return m.getExpressionsFunc(token)
}

func (m *mockTaskService) SplitTasks(expression *models.Expression) {
}

func (m *mockTaskService) Register(req *models.RegisterRequest) error {
	return m.registerFunc(req)
}

func (m *mockTaskService) Login(req *models.Login) (string, error) {
	return m.loginFunc(req)
}

func TestRegisterHandler(t *testing.T) {
	tests := []struct {
		name           string
		request        models.RegisterRequest
		mockRegister   func(*models.RegisterRequest) error
		expectedStatus int
	}{
		{
			name: "successful registration",
			request: models.RegisterRequest{
				Username: "testuser",
				Password: "testpass",
			},
			mockRegister: func(req *models.RegisterRequest) error {
				return nil
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "duplicate username",
			request: models.RegisterRequest{
				Username: "existinguser",
				Password: "testpass",
			},
			mockRegister: func(req *models.RegisterRequest) error {
				return errors.New("username already exists")
			},
			expectedStatus: http.StatusUnprocessableEntity,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &mockTaskService{
				registerFunc: tt.mockRegister,
			}

			orch := &http2.Orchestrator{
				Service: service,
				Ctx:     context.Background(),
				Config:  &config.Config{},
			}

			body, _ := json.Marshal(tt.request)

			req := httptest.NewRequest("POST", "/api/v1/register", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()

			handler := http2.RegisterHandler(orch)
			handler.ServeHTTP(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rr.Code)
			}
		})
	}
}

func TestCalcHandler(t *testing.T) {
	tests := []struct {
		name           string
		request        models.Request
		token          string
		mockCreateExpr func(string, string) (*models.Expression, error)
		expectedStatus int
	}{
		{
			name: "successful expression creation",
			request: models.Request{
				Expression: "2+2",
			},
			token: "valid-token",
			mockCreateExpr: func(expr, token string) (*models.Expression, error) {
				return &models.Expression{Id: "123", Status: "created", Result: 4}, nil
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "invalid expression",
			request: models.Request{
				Expression: "invalid",
			},
			token: "valid-token",
			mockCreateExpr: func(expr, token string) (*models.Expression, error) {
				return &models.Expression{}, errors.New("invalid expression")
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "unauthorized",
			request: models.Request{
				Expression: "2+2",
			},
			token:          "",
			mockCreateExpr: nil,
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &mockTaskService{
				createExprFunc: tt.mockCreateExpr,
			}

			orch := &http2.Orchestrator{
				Service: service,
				Ctx:     context.Background(),
				Config:  &config.Config{},
			}

			body, _ := json.Marshal(tt.request)
			req := httptest.NewRequest("POST", "/api/v1/calculate", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			if tt.token != "" {
				req.Header.Set("Authorization", "Bearer "+tt.token)
			}

			rr := httptest.NewRecorder()
			handler := http2.CalcHandler(orch)
			handler.ServeHTTP(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rr.Code)
			}
		})
	}
}

func TestExpressionHandler(t *testing.T) {
	tests := []struct {
		name           string
		id             string
		token          string
		mockGetExpr    func(string, string) (*models.Expression, error)
		expectedStatus int
	}{
		{
			name:  "successful get",
			id:    "123",
			token: "valid-token",
			mockGetExpr: func(id, token string) (*models.Expression, error) {
				return &models.Expression{Id: id, Status: "completed", Result: 4}, nil
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:  "not found",
			id:    "456",
			token: "valid-token",
			mockGetExpr: func(id, token string) (*models.Expression, error) {
				return &models.Expression{}, errors.New("not found")
			},
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &mockTaskService{
				getExpressionFunc: tt.mockGetExpr,
			}

			orch := &http2.Orchestrator{
				Service: service,
				Ctx:     context.Background(),
				Config:  &config.Config{},
			}

			req := httptest.NewRequest("GET", "/api/v1/expressions/"+tt.id, nil)
			req.Header.Set("Authorization", "Bearer "+tt.token)

			rr := httptest.NewRecorder()

			router := mux.NewRouter()
			router.HandleFunc("/api/v1/expressions/{id}", http2.ExpressionHandler(orch))
			router.ServeHTTP(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rr.Code)
			}
		})
	}
}
