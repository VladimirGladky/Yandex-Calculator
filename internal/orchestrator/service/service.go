package service

import (
	"fmt"
	"github.com/VladimirGladky/FinalTaskFirstSprint/internal/models"
	"github.com/VladimirGladky/FinalTaskFirstSprint/internal/orchestrator/parser"
	"github.com/google/uuid"
	"os"
	"strconv"
	"sync"
)

type Service struct {
	Mu             sync.Mutex
	ExpressionsMap map[string]*models.Expression
	TasksMap       map[string]*models.Task
	TasksArr       []*models.Task
	TaskCounter    int
}

func NewService() *Service {
	return &Service{
		ExpressionsMap: make(map[string]*models.Expression),
		TasksMap:       make(map[string]*models.Task),
		TasksArr:       make([]*models.Task, 0),
	}
}

func (s *Service) Login(lp *models.Login) (string, error) {
	if err := lp.Validate(); err != nil {
		return "", err
	}
	return "", fmt.Errorf("wrong login or password")
}

func (s *Service) Register(rp *models.RegisterRequest) error {
	if err := rp.Validate(); err != nil {
		return err
	}
	return fmt.Errorf("user already exists")
}

func (s *Service) GetExpression(id string) (*models.Expression, error) {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	if _, ok := s.ExpressionsMap[id]; !ok {
		return nil, fmt.Errorf("expression with id %s not found", id)
	}
	return s.ExpressionsMap[id], nil
}

func (s *Service) CreateExpression(expr string) (*models.Expression, error) {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	ast, err := parser.BuildExpressionTree(expr)
	if err != nil {
		return nil, err
	}

	id := uuid.New().String()
	expression := &models.Expression{
		Id:     id,
		Status: "in progress",
		Ast:    ast,
	}
	s.ExpressionsMap[id] = expression
	s.SplitTasks(s.ExpressionsMap[id])
	return expression, nil
}

func (s *Service) GetExpressions() []models.Expression {
	expressions := make([]models.Expression, 0, len(s.ExpressionsMap))
	for _, v := range s.ExpressionsMap {
		if v.Ast != nil && v.Ast.IsLeaf {
			v.Status = "done"
			v.Result = v.Ast.Value
		}
		expressions = append(expressions, *v)
	}
	return expressions
}

func (s *Service) SplitTasks(expression *models.Expression) {
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

			task := &models.Task{
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
