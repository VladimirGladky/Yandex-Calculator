package taskgRPC

import (
	"context"
	task2 "github.com/VladimirGladky/FinalTaskFirstSprint/gen/proto/task"
	"github.com/VladimirGladky/FinalTaskFirstSprint/internal/orchestrator/transport/http"
	"github.com/VladimirGladky/FinalTaskFirstSprint/pkg/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type Service struct {
	task2.UnimplementedTaskManagementServiceServer
	orchestrator *http.Orchestrator
}

func NewService(orchestrator *http.Orchestrator) *Service {
	return &Service{
		orchestrator: orchestrator,
	}
}

func Register(s *grpc.Server, orchestrator *http.Orchestrator) {
	task2.RegisterTaskManagementServiceServer(s, NewService(orchestrator))
}

func (s *Service) TaskGet(context.Context, *task2.TaskGetRequest) (*task2.TaskGetResponse, error) {
	s.orchestrator.Service.Mu.Lock()
	defer s.orchestrator.Service.Mu.Unlock()
	if len(s.orchestrator.Service.TasksArr) == 0 {
		return nil, nil
	}
	taskGet := s.orchestrator.Service.TasksArr[0]
	s.orchestrator.Service.TasksArr = s.orchestrator.Service.TasksArr[1:]
	if expr, exists := s.orchestrator.Service.ExpressionsMap[taskGet.ExprID]; exists {
		expr.Status = "in_progress"
	}
	logger.GetLoggerFromCtx(s.orchestrator.Ctx).Info("TaskGet", zap.Any("taskGet", &taskGet))
	return &task2.TaskGetResponse{Id: taskGet.ID, Arg1: float32(taskGet.Arg1), Arg2: float32(taskGet.Arg2), Operation: taskGet.Operation, OperationTime: int32(taskGet.OperationTime)}, nil
}

func (s *Service) TaskPost(ctx context.Context, req *task2.TaskPostRequest) (*task2.TaskPostResponse, error) {
	s.orchestrator.Service.Mu.Lock()
	taskPost, ok := s.orchestrator.Service.TasksMap[req.Id]
	if !ok {
		s.orchestrator.Service.Mu.Unlock()
		logger.GetLoggerFromCtx(s.orchestrator.Ctx).Info("No taskPost available")
		return nil, nil
	}
	taskPost.Node.IsLeaf, taskPost.Node.Value = true, float64(req.Result)
	delete(s.orchestrator.Service.TasksMap, req.Id)
	if expression, ex := s.orchestrator.Service.ExpressionsMap[taskPost.ExprID]; ex {
		s.orchestrator.Service.SplitTasks(expression)
		if expression.Ast.IsLeaf {
			expression.Status, expression.Result = "completed", expression.Ast.Value
		}
	}
	logger.GetLoggerFromCtx(s.orchestrator.Ctx).Info("TaskPost", zap.Any("taskPost", &taskPost))
	s.orchestrator.Service.Mu.Unlock()
	return &task2.TaskPostResponse{}, nil
}
