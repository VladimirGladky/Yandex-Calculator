package models

import "github.com/VladimirGladky/FinalTaskFirstSprint/internal/orchestrator/parser"

type TaskGet struct {
	Id            string  `json:"id"`
	Arg1          float64 `json:"arg1"`
	Arg2          float64 `json:"arg2"`
	Operation     string  `json:"operation"`
	OperationTime int     `json:"operation_time"`
}

type Task struct {
	ID            string                 `json:"id"`
	Arg1          float64                `json:"arg1"`
	Arg2          float64                `json:"arg2"`
	Operation     string                 `json:"operation"`
	OperationTime int                    `json:"operation_time"`
	Node          *parser.ExpressionNode `json:"-"`
	ExprID        string                 `json:"-"`
}
