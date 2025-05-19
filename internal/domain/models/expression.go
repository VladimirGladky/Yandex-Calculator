package models

import "github.com/VladimirGladky/Yandex-Calculator/internal/orchestrator/parser"

type Expression struct {
	Id     string                 `json:"id"`
	Status string                 `json:"status"`
	Result float64                `json:"result"`
	Ast    *parser.ExpressionNode `json:"-"`
}
