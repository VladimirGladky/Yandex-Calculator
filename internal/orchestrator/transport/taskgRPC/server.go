package taskgRPC

import (
	"context"
	"fmt"
	task2 "github.com/VladimirGladky/FinalTaskFirstSprint/gen/proto/task"
	"github.com/VladimirGladky/FinalTaskFirstSprint/internal/orchestrator/service"
	"google.golang.org/grpc"
)

type Service struct {
	task2.UnimplementedTaskManagementServiceServer
	taskService *service.Service
}

func NewService(taskService *service.Service) *Service {
	return &Service{
		taskService: taskService,
	}
}

func Register(s *grpc.Server, taskService *service.Service) {
	task2.RegisterTaskManagementServiceServer(s, NewService(taskService))
}

func (s *Service) TaskGet(context.Context, *task2.TaskGetRequest) (*task2.TaskGetResponse, error) {
	s.taskService.Mu.Lock()
	defer s.taskService.Mu.Unlock()
	if len(s.taskService.TasksArr) == 0 {
		return nil, nil
	}
	taskGet := s.taskService.TasksArr[0]
	s.taskService.TasksArr = s.taskService.TasksArr[1:]
	if expr, exists := s.taskService.ExpressionsMap[taskGet.ExprID]; exists {
		expr.Status = "in_progress"
	}
	return &task2.TaskGetResponse{Id: taskGet.ID, Arg1: float32(taskGet.Arg1), Arg2: float32(taskGet.Arg2), Operation: taskGet.Operation, OperationTime: int32(taskGet.OperationTime)}, nil
}

func (s *Service) TaskPost(ctx context.Context, req *task2.TaskPostRequest) (*task2.TaskPostResponse, error) {
	s.taskService.Mu.Lock()
	taskPost, ok := s.taskService.TasksMap[req.Id]
	if !ok {
		s.taskService.Mu.Unlock()
		return nil, fmt.Errorf("task with id %s not found", req.Id)
	}
	ctx = context.WithValue(ctx, "email", req.Result)
	taskPost.Node.IsLeaf, taskPost.Node.Value = true, float64(req.Result)
	delete(s.taskService.TasksMap, req.Id)
	if expression, ex := s.taskService.ExpressionsMap[taskPost.ExprID]; ex {
		s.taskService.SplitTasks(expression)
		if expression.Ast.IsLeaf {
			expression.Status, expression.Result = "completed", expression.Ast.Value
			err := s.taskService.UpdateExpression(expression)
			if err != nil {
				return nil, err
			}
		}
	}
	s.taskService.Mu.Unlock()
	return &task2.TaskPostResponse{}, nil
}
