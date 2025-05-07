package calculation

import (
	"github.com/VladimirGladky/FinalTaskFirstSprint/internal/domain/models"
	"time"
)

func ComputeTask(task models.TaskGet) (float64, error) {
	switch task.Operation {
	case "+":
		time.Sleep(time.Duration(task.OperationTime) * time.Millisecond)
		return task.Arg1 + task.Arg2, nil
	case "-":
		time.Sleep(time.Duration(task.OperationTime) * time.Millisecond)
		return task.Arg1 - task.Arg2, nil
	case "*":
		time.Sleep(time.Duration(task.OperationTime) * time.Millisecond)
		return task.Arg1 * task.Arg2, nil
	case "/":
		time.Sleep(time.Duration(task.OperationTime) * time.Millisecond)
		if task.Arg2 == 0 {
			return 0, ErrDivisionByZero
		}
		return task.Arg1 / task.Arg2, nil
	default:
		return 0, ErrInvalidOperation
	}
}
