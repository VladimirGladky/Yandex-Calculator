package parser

import "fmt"

var (
	EmptyExpression   = fmt.Errorf("пустое выражение")
	InvalidExpression = fmt.Errorf("неверное выражение")
	NoBracket         = fmt.Errorf("нет скобки")
	ExpectedNumber    = fmt.Errorf("ожидалось число")
	InvalidNumber     = fmt.Errorf("неверное число")
)
