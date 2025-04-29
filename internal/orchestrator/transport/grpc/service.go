package grpc

import (
	"context"
	task2 "github.com/VladimirGladky/FinalTaskFirstSprint/gen/proto/task"
	"github.com/VladimirGladky/FinalTaskFirstSprint/internal/orchestrator/server"
	"github.com/VladimirGladky/FinalTaskFirstSprint/pkg/logger"
	"go.uber.org/zap"
)

type Service struct {
	task2.UnimplementedTaskManagementServiceServer
	orchestrator *server.Orchestrator
}

func NewService(orchestrator *server.Orchestrator) *Service {
	return &Service{
		orchestrator: orchestrator,
	}
}

func (s *Service) TaskGet(context.Context, *task2.TaskGetRequest) (*task2.TaskGetResponse, error) {
	s.orchestrator.Mu.Lock()
	defer s.orchestrator.Mu.Unlock()
	if len(s.orchestrator.TasksArr) == 0 {
		return nil, nil
	}
	taskGet := s.orchestrator.TasksArr[0]
	s.orchestrator.TasksArr = s.orchestrator.TasksArr[1:]
	if expr, exists := s.orchestrator.ExpressionsMap[taskGet.ExprID]; exists {
		expr.Status = "in_progress"
	}
	logger.GetLoggerFromCtx(s.orchestrator.Ctx).Info(s.orchestrator.Ctx, "TaskGet", zap.Any("taskGet", &taskGet))
	return &task2.TaskGetResponse{Id: taskGet.ID, Arg1: float32(taskGet.Arg1), Arg2: float32(taskGet.Arg2), Operation: taskGet.Operation, OperationTime: int32(taskGet.OperationTime)}, nil
}

func (s *Service) TaskPost(ctx context.Context, req *task2.TaskPostRequest) (*task2.TaskPostResponse, error) {
	s.orchestrator.Mu.Lock()
	taskPost, ok := s.orchestrator.TasksMap[req.Id]
	if !ok {
		s.orchestrator.Mu.Unlock()
		logger.GetLoggerFromCtx(s.orchestrator.Ctx).Info(s.orchestrator.Ctx, "No taskPost available")
		return nil, nil
	}
	taskPost.Node.IsLeaf, taskPost.Node.Value = true, float64(req.Result)
	delete(s.orchestrator.TasksMap, req.Id)
	if expression, ex := s.orchestrator.ExpressionsMap[taskPost.ExprID]; ex {
		s.orchestrator.SplitTasks(expression)
		if expression.Ast.IsLeaf {
			expression.Status, expression.Result = "completed", expression.Ast.Value
		}
	}
	logger.GetLoggerFromCtx(s.orchestrator.Ctx).Info(s.orchestrator.Ctx, "TaskPost", zap.Any("taskPost", &taskPost))
	s.orchestrator.Mu.Unlock()
	return &task2.TaskPostResponse{}, nil
}
