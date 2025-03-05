package calculation

import "errors"

var (
	ErrInvalidOperation = errors.New("invalid operation")
	ErrDivisionByZero   = errors.New("division by zero")
)
