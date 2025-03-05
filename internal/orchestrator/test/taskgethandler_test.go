package test

import (
	"context"
	"encoding/json"
	"github.com/VladimirGladky/FinalTaskFirstSprint/internal/models"
	"github.com/VladimirGladky/FinalTaskFirstSprint/internal/orchestrator/server"
	"github.com/VladimirGladky/FinalTaskFirstSprint/pkg/logger"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTaskHandlerGet(t *testing.T) {
	ctx := context.Background()
	ctx, _ = logger.New(ctx)
	orchestrator := server.New(ctx)

	router := mux.NewRouter()
	router.HandleFunc("/internal/task", server.TaskHandlerGet(orchestrator)).Methods("GET")

	t.Run("No tasks available", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/internal/task", nil)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusNotFound, rr.Code)
		assert.JSONEq(t, `{"error":"No task available"}`, rr.Body.String())
	})

	t.Run("Task successfully retrieved", func(t *testing.T) {
		task := &models.Task{
			ID:            "1",
			ExprID:        "expr1",
			Arg1:          10,
			Arg2:          5,
			Operation:     "+",
			OperationTime: 100,
		}
		orchestrator.TasksArr = append(orchestrator.TasksArr, task)
		orchestrator.ExpressionsMap["expr1"] = &models.Expression{Id: "expr1", Status: "pending"}

		req, err := http.NewRequest("GET", "/internal/task", nil)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var responseTask models.Task
		err = json.NewDecoder(rr.Body).Decode(&responseTask)
		assert.NoError(t, err)
		assert.Equal(t, task.ID, responseTask.ID)
		assert.Equal(t, task.Arg1, responseTask.Arg1)
		assert.Equal(t, task.Arg2, responseTask.Arg2)
		assert.Equal(t, task.Operation, responseTask.Operation)

		assert.Empty(t, orchestrator.TasksArr)
	})
}
