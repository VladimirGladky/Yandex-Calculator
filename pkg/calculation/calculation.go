package calculation

import (
	"strconv"
	"strings"
	"unicode"
)

func Calc(expression string) (float64, error) {
	var numbersStack []float64
	var operationsStack []rune
	var i int
	var leftBrackets int
	var rightBrackets int
	var err error
	expression = strings.ReplaceAll(expression, " ", "")

	if !isValidExpression(expression) {
		return 0, ErrInvalidExpression
	}

	for i < len(expression) {
		if unicode.IsDigit(rune(expression[i])) {
			num, _ := strconv.ParseFloat(string(expression[i]), 64)
			numbersStack = append(numbersStack, num)
			i++
			continue
		}
		switch expression[i] {
		case '+', '-', '*', '/':
			for len(operationsStack) > 0 && priorityOperations(operationsStack[len(operationsStack)-1]) >= priorityOperations(rune(expression[i])) {
				numbersStack, operationsStack, err = makeOperation(numbersStack, operationsStack)
				if err != nil {
					return 0, err
				}
			}
			operationsStack = append(operationsStack, rune(expression[i]))
			i++
			continue
		}

		if expression[i] == '(' {
			operationsStack = append(operationsStack, rune(expression[i]))
			leftBrackets++
			i++
			continue
		}

		if expression[i] == ')' {
			rightBrackets++
			if leftBrackets >= rightBrackets {
				for operationsStack[len(operationsStack)-1] != '(' {
					numbersStack, operationsStack, err = makeOperation(numbersStack, operationsStack)
					if err != nil {
						return 0, err
					}
				}
				operationsStack = operationsStack[:len(operationsStack)-1]
				leftBrackets--
				rightBrackets--
			} else {
				return 0, ErrInvalidExpression
			}
			i++
			continue
		}
	}

	if leftBrackets != 0 || rightBrackets != 0 {
		return 0, ErrInvalidExpression
	}

	if len(numbersStack)-1 != len(operationsStack) {
		return 0, ErrInvalidExpression
	}

	for len(operationsStack) > 0 {
		numbersStack, operationsStack, err = makeOperation(numbersStack, operationsStack)
		if err != nil {
			return 0, err
		}
	}

	return numbersStack[0], nil
}

func makeOperation(numbersStack []float64, operationsStack []rune) ([]float64, []rune, error) {
	if len(numbersStack) < 2 || len(operationsStack) == 0 {
		return numbersStack, operationsStack, nil
	}
	a := numbersStack[len(numbersStack)-2]
	b := numbersStack[len(numbersStack)-1]
	operation := operationsStack[len(operationsStack)-1]

	var result float64
	switch operation {
	case '+':
		result = a + b
	case '-':
		result = a - b
	case '*':
		result = a * b
	case '/':
		if b == 0 {
			return numbersStack, operationsStack, ErrDivisionByZero
		}
		result = a / b
	default:
		return numbersStack, operationsStack, ErrInvalidExpression
	}

	numbersStack = numbersStack[:len(numbersStack)-2]
	operationsStack = operationsStack[:len(operationsStack)-1]
	return append(numbersStack, result), operationsStack, nil
}

func priorityOperations(operation rune) int {
	switch operation {
	case '*', '/':
		return 2
	case '+', '-':
		return 1
	}
	return 0
}

func isValidExpression(expression string) bool {
	for _, char := range expression {
		if !strings.ContainsRune("0123456789+-*/()", char) {
			return false
		}
	}

	var countNumbers int
	var countOperations int
	for _, char := range expression {
		if strings.ContainsRune("0123456789", char) {
			countNumbers++
		} else if strings.ContainsRune("+-*/", char) {
			countOperations++
		}
	}

	return countNumbers > countOperations
}
