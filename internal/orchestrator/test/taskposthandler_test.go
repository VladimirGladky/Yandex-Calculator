package test

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/VladimirGladky/FinalTaskFirstSprint/internal/models"
	"github.com/VladimirGladky/FinalTaskFirstSprint/internal/orchestrator/parser"
	"github.com/VladimirGladky/FinalTaskFirstSprint/internal/orchestrator/server"
	"github.com/VladimirGladky/FinalTaskFirstSprint/pkg/logger"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTaskHandlerPost(t *testing.T) {
	ctx := context.Background()
	ctx, _ = logger.New(ctx)
	orchestrator := server.New(ctx)

	router := mux.NewRouter()
	router.HandleFunc("/internal/task", server.TaskHandlerPost(orchestrator)).Methods("POST")

	t.Run("Task successfully completed", func(t *testing.T) {
		task := &models.Task{
			ID:            "1",
			ExprID:        "expr1",
			Arg1:          10,
			Arg2:          5,
			Operation:     "+",
			OperationTime: 100,
			Node: &parser.ExpressionNode{
				IsLeaf: false,
				Value:  0,
			},
		}
		orchestrator.TasksMap["1"] = task
		orchestrator.ExpressionsMap["expr1"] = &models.Expression{
			Id:     "expr1",
			Status: "pending",
			Ast:    task.Node,
		}

		requestBody := models.TaskPost{
			Id:     "1",
			Result: 15,
		}
		requestBodyBytes, err := json.Marshal(requestBody)
		assert.NoError(t, err)

		req, err := http.NewRequest("POST", "/internal/task", bytes.NewBuffer(requestBodyBytes))
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		_, exists := orchestrator.TasksMap["1"]
		assert.False(t, exists)

		expr, exists := orchestrator.ExpressionsMap["expr1"]
		assert.True(t, exists)
		assert.Equal(t, "completed", expr.Status)
		assert.Equal(t, 15.0, expr.Result)
	})

	t.Run("Task not found", func(t *testing.T) {
		requestBody := models.TaskPost{
			Id:     "999",
			Result: 15,
		}
		requestBodyBytes, err := json.Marshal(requestBody)
		assert.NoError(t, err)

		req, err := http.NewRequest("POST", "/internal/task", bytes.NewBuffer(requestBodyBytes))
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusNotFound, rr.Code)
		assert.JSONEq(t, `{"error":"Task not found"}`, rr.Body.String())
	})
}
