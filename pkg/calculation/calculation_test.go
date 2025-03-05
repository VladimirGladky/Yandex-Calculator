package calculation

import (
	"testing"
	"time"

	"github.com/VladimirGladky/FinalTaskFirstSprint/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestComputeTask(t *testing.T) {
	t.Run("Addition", func(t *testing.T) {
		task := models.TaskGet{
			Arg1:          10,
			Arg2:          5,
			Operation:     "+",
			OperationTime: 100,
		}

		result, err := ComputeTask(task)
		assert.NoError(t, err)
		assert.Equal(t, 15.0, result)
	})

	t.Run("Subtraction", func(t *testing.T) {
		task := models.TaskGet{
			Arg1:          10,
			Arg2:          5,
			Operation:     "-",
			OperationTime: 100,
		}

		result, err := ComputeTask(task)
		assert.NoError(t, err)
		assert.Equal(t, 5.0, result)
	})

	t.Run("Multiplication", func(t *testing.T) {
		task := models.TaskGet{
			Arg1:          10,
			Arg2:          5,
			Operation:     "*",
			OperationTime: 100,
		}

		result, err := ComputeTask(task)
		assert.NoError(t, err)
		assert.Equal(t, 50.0, result)
	})

	t.Run("Division", func(t *testing.T) {
		task := models.TaskGet{
			Arg1:          10,
			Arg2:          5,
			Operation:     "/",
			OperationTime: 100,
		}

		result, err := ComputeTask(task)
		assert.NoError(t, err)
		assert.Equal(t, 2.0, result)
	})

	t.Run("Division by zero", func(t *testing.T) {
		task := models.TaskGet{
			Arg1:          10,
			Arg2:          0,
			Operation:     "/",
			OperationTime: 100,
		}

		result, err := ComputeTask(task)
		assert.Error(t, err)
		assert.Equal(t, ErrDivisionByZero, err)
		assert.Equal(t, 0.0, result)
	})

	t.Run("Invalid operation", func(t *testing.T) {
		task := models.TaskGet{
			Arg1:          10,
			Arg2:          5,
			Operation:     "%",
			OperationTime: 100,
		}

		result, err := ComputeTask(task)
		assert.Error(t, err)
		assert.Equal(t, ErrInvalidOperation, err)
		assert.Equal(t, 0.0, result)
	})

	t.Run("Operation time", func(t *testing.T) {
		task := models.TaskGet{
			Arg1:          10,
			Arg2:          5,
			Operation:     "+",
			OperationTime: 100,
		}

		start := time.Now()
		_, err := ComputeTask(task)
		duration := time.Since(start)

		assert.NoError(t, err)
		assert.True(t, duration >= 100*time.Millisecond, "Operation should take at least 100ms")
	})
}
