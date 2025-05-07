package service

import (
	"context"
	"fmt"
	models2 "github.com/VladimirGladky/FinalTaskFirstSprint/internal/domain/models"
	"github.com/VladimirGladky/FinalTaskFirstSprint/internal/orchestrator/parser"
	"github.com/VladimirGladky/FinalTaskFirstSprint/internal/orchestrator/storage/sqlite"
	"github.com/VladimirGladky/FinalTaskFirstSprint/pkg/jwt"
	"github.com/VladimirGladky/FinalTaskFirstSprint/pkg/logger"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"os"
	"strconv"
	"sync"
	"time"
)

type TaskServiceInterface interface {
	Login(lp *models2.Login) (string, error)
	Register(rp *models2.RegisterRequest) error
	GetExpression(id string) (*models2.Expression, error)
	CreateExpression(expr string) (*models2.Expression, error)
	GetExpressions() []models2.Expression
	SplitTasks(expression *models2.Expression)
}

type Service struct {
	storage        *sqlite.Storage
	Mu             sync.Mutex
	ExpressionsMap map[string]*models2.Expression
	TasksMap       map[string]*models2.Task
	TasksArr       []*models2.Task
	TaskCounter    int
	ctx            context.Context
}

func NewService(ctx context.Context, storage *sqlite.Storage) *Service {
	return &Service{
		ExpressionsMap: make(map[string]*models2.Expression),
		TasksMap:       make(map[string]*models2.Task),
		TasksArr:       make([]*models2.Task, 0),
		ctx:            ctx,
		storage:        storage,
	}
}

func (s *Service) Login(lp *models2.Login) (string, error) {
	if err := lp.Validate(); err != nil {
		logger.GetLoggerFromCtx(s.ctx).Error("error validating login: %v", zap.Error(err))
		return "", err
	}
	user, err := s.storage.User(s.ctx, lp.Username)
	if err != nil {
		logger.GetLoggerFromCtx(s.ctx).Error("error getting user: %v", zap.Error(err))
		return "", fmt.Errorf("error getting user: %v", err)
	}
	err = bcrypt.CompareHashAndPassword(user.PassHash, []byte(lp.Password))
	if err != nil {
		logger.GetLoggerFromCtx(s.ctx).Error("wrong login or password")
		return "", fmt.Errorf("wrong login or password")
	}
	token, err := jwt.NewToken(user, 24*time.Hour)
	if err != nil {
		logger.GetLoggerFromCtx(s.ctx).Error("error creating token: %v", zap.Error(err))
		return "", fmt.Errorf("error creating token: %v", err)
	}
	return token, nil
}

func (s *Service) Register(rp *models2.RegisterRequest) error {
	if err := rp.Validate(); err != nil {
		logger.GetLoggerFromCtx(s.ctx).Error("error validating register request: %v", zap.Error(err))
		return fmt.Errorf("error validating register request: %v", err)
	}
	passHash, err := bcrypt.GenerateFromPassword([]byte(rp.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.GetLoggerFromCtx(s.ctx).Error("error hashing password: %v", zap.Error(err))
		return fmt.Errorf("error hashing password: %v", err)
	}
	err = s.storage.SaveUser(s.ctx, &models2.User{
		Email:    rp.Username,
		PassHash: passHash,
	})
	if err != nil {
		logger.GetLoggerFromCtx(s.ctx).Error("error saving user: %v", zap.Error(err))
		return fmt.Errorf("error saving user: %v", err)
	}
	logger.GetLoggerFromCtx(s.ctx).Info("user registered")
	return nil
}

func (s *Service) GetExpression(id string, token string) (*models2.Expression, error) {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	_, err := jwt.GetEmailFromToken(token)
	if err != nil {
		return nil, err
	}
	if _, ok := s.ExpressionsMap[id]; !ok {
		return nil, fmt.Errorf("expression with id %s not found", id)
	}
	return s.ExpressionsMap[id], nil
}

func (s *Service) CreateExpression(expr string, token string) (*models2.Expression, error) {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	_, err := jwt.GetEmailFromToken(token)
	if err != nil {
		return nil, err
	}

	ast, err := parser.BuildExpressionTree(expr)
	if err != nil {
		return nil, err
	}

	id := uuid.New().String()
	expression := &models2.Expression{
		Id:     id,
		Status: "in progress",
		Ast:    ast,
	}
	s.ExpressionsMap[id] = expression
	s.SplitTasks(s.ExpressionsMap[id])
	return expression, nil
}

func (s *Service) GetExpressions(token string) ([]models2.Expression, error) {
	_, err := jwt.GetEmailFromToken(token)
	if err != nil {
		return nil, err
	}
	expressions := make([]models2.Expression, 0, len(s.ExpressionsMap))
	for _, v := range s.ExpressionsMap {
		if v.Ast != nil && v.Ast.IsLeaf {
			v.Status = "done"
			v.Result = v.Ast.Value
		}
		expressions = append(expressions, *v)
	}
	return expressions, nil
}

func (s *Service) SplitTasks(expression *models2.Expression) {
	var visitNode func(node *parser.ExpressionNode)
	visitNode = func(node *parser.ExpressionNode) {
		if node == nil || node.IsLeaf {
			return
		}

		visitNode(node.Left)
		visitNode(node.Right)

		if node.Left != nil && node.Right != nil && node.Left.IsLeaf && node.Right.IsLeaf && !node.TaskScheduled {
			s.TaskCounter++
			taskID := fmt.Sprintf("%d", s.TaskCounter)
			opTime := getOperationTime(node.Operator)

			task := &models2.Task{
				ID:            taskID,
				ExprID:        expression.Id,
				Arg1:          node.Left.Value,
				Arg2:          node.Right.Value,
				Operation:     node.Operator,
				OperationTime: opTime,
				Node:          node,
			}

			node.TaskScheduled = true
			s.TasksMap[taskID] = task
			s.TasksArr = append(s.TasksArr, task)
		}
	}

	visitNode(expression.Ast)
}

func getOperationTime(operator string) int {
	var time int
	var err error

	switch operator {
	case "+":
		time, err = strconv.Atoi(os.Getenv("TIME_ADDITION_MS"))
	case "-":
		time, err = strconv.Atoi(os.Getenv("TIME_SUBTRACTION_MS"))
	case "*":
		time, err = strconv.Atoi(os.Getenv("TIME_MULTIPLICATIONS_MS"))
	case "/":
		time, err = strconv.Atoi(os.Getenv("TIME_DIVISIONS_MS"))
	default:
		time = 100
	}

	if err != nil {
		time = 100
	}

	return time
}
