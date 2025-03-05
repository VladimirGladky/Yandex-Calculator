package parser

import (
	"strconv"
	"unicode"
)

type ExpressionParser struct {
	input  string
	cursor int
}

func (ep *ExpressionParser) nextChar() rune {
	if ep.cursor < len(ep.input) {
		return rune(ep.input[ep.cursor])
	}
	return 0
}

func (ep *ExpressionParser) consumeChar() rune {
	char := ep.nextChar()
	ep.cursor += 1
	return char
}

func (ep *ExpressionParser) parse() (*ExpressionNode, error) {
	node, err := ep.parseAdditionOrSubtraction()
	if err != nil {
		return nil, err
	}
	return node, nil
}

func (ep *ExpressionParser) parseAdditionOrSubtraction() (*ExpressionNode, error) {
	node, err := ep.parseMultiplicationOrDivision()
	if err != nil {
		return nil, err
	}
	for {
		char := ep.nextChar()
		if char == '+' || char == '-' {
			op := string(ep.consumeChar())
			rightNode, err := ep.parseMultiplicationOrDivision()
			if err != nil {
				return nil, err
			}
			node = &ExpressionNode{
				IsLeaf:   false,
				Operator: op,
				Left:     node,
				Right:    rightNode,
			}
		} else {
			break
		}
	}
	return node, nil
}

func (ep *ExpressionParser) parseMultiplicationOrDivision() (*ExpressionNode, error) {
	node, err := ep.parsePrimary()
	if err != nil {
		return nil, err
	}
	for {
		char := ep.nextChar()
		if char == '*' || char == '/' {
			op := string(ep.consumeChar())
			rightNode, err := ep.parsePrimary()
			if err != nil {
				return nil, err
			}
			node = &ExpressionNode{
				IsLeaf:   false,
				Operator: op,
				Left:     node,
				Right:    rightNode,
			}
		} else {
			break
		}
	}
	return node, nil
}

func (ep *ExpressionParser) parsePrimary() (*ExpressionNode, error) {
	char := ep.nextChar()
	if char == '(' {
		ep.consumeChar()
		node, err := ep.parse()
		if err != nil {
			return nil, err
		}
		if ep.nextChar() != ')' {
			return nil, NoBracket
		}
		ep.consumeChar()
		return node, nil
	}
	startPos := ep.cursor
	if char == '-' {
		ep.consumeChar()
	}
	for {
		char = ep.nextChar()
		if unicode.IsDigit(char) || char == '.' {
			ep.consumeChar()
		} else {
			break
		}
	}
	numStr := ep.input[startPos:ep.cursor]
	if numStr == "" {
		return nil, ExpectedNumber
	}
	num, err := strconv.ParseFloat(numStr, 64)
	if err != nil {
		return nil, InvalidNumber
	}
	return &ExpressionNode{
		IsLeaf: true,
		Value:  num,
	}, nil
}
