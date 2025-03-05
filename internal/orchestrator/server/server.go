package server

import (
	"context"
	"fmt"
	"github.com/VladimirGladky/FinalTaskFirstSprint/internal/models"
	"github.com/VladimirGladky/FinalTaskFirstSprint/internal/orchestrator/parser"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"strconv"
	"sync"
)

type Config struct {
	Addr string
}

func ConfigFromEnv() *Config {
	config := new(Config)
	config.Addr = os.Getenv("PORT")
	if config.Addr == "" {
		config.Addr = "4040"
	}
	return config
}

type Orchestrator struct {
	ctx            context.Context
	config         *Config
	mu             sync.Mutex
	ExpressionsMap map[string]*models.Expression
	TasksArr       []*models.Task
	TasksMap       map[string]*models.Task
	TaskCounter    int
}

func New(ctx context.Context) *Orchestrator {
	return &Orchestrator{
		ctx:            ctx,
		config:         ConfigFromEnv(),
		ExpressionsMap: make(map[string]*models.Expression),
		TasksMap:       make(map[string]*models.Task),
		TasksArr:       make([]*models.Task, 0),
	}
}

func (o *Orchestrator) Run() error {
	router := mux.NewRouter()
	router.HandleFunc("/api/v1/calculate", CalcHandler(o))
	router.HandleFunc("/api/v1/expressions", ExpressionsHandler(o))
	router.HandleFunc("/api/v1/expressions/{id}", ExpressionHandler(o))
	router.HandleFunc("/internal/task", TaskHandlerPost(o)).Methods("POST")
	router.HandleFunc("/internal/task", TaskHandlerGet(o)).Methods("GET")
	return http.ListenAndServe(":"+o.config.Addr, router)
}

func (o *Orchestrator) SplitTasks(expr *models.Expression) {
	var visitNode func(node *parser.ExpressionNode)
	visitNode = func(node *parser.ExpressionNode) {
		if node == nil || node.IsLeaf {
			return
		}

		visitNode(node.Left)
		visitNode(node.Right)

		if node.Left != nil && node.Right != nil && node.Left.IsLeaf && node.Right.IsLeaf && !node.TaskScheduled {
			o.TaskCounter++
			taskID := fmt.Sprintf("%d", o.TaskCounter)
			opTime := getOperationTime(node.Operator)

			task := &models.Task{
				ID:            taskID,
				ExprID:        expr.Id,
				Arg1:          node.Left.Value,
				Arg2:          node.Right.Value,
				Operation:     node.Operator,
				OperationTime: opTime,
				Node:          node,
			}

			node.TaskScheduled = true
			o.TasksMap[taskID] = task
			o.TasksArr = append(o.TasksArr, task)
		}
	}

	visitNode(expr.Ast)
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
