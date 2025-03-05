package parser

import (
	"strings"
)

func BuildExpressionTree(expression string) (*ExpressionNode, error) {
	cleanedExpr := strings.ReplaceAll(expression, " ", "")
	if cleanedExpr == "" {
		return nil, EmptyExpression
	}
	p := &ExpressionParser{input: cleanedExpr, cursor: 0}
	node, err := p.parse()
	if err != nil {
		return nil, err
	}
	if p.cursor < len(p.input) {
		return nil, InvalidExpression
	}
	return node, nil
}
