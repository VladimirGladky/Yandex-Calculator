package calculation_test

import (
	"github.com/VladimirGladky/FinalTaskFirstSprint/pkg/calculation"
	"testing"
)

func TestCalc(t *testing.T) {
	testCasesSuccess := []struct {
		name           string
		expression     string
		expectedResult float64
	}{
		{
			name:           "simple",
			expression:     "1+1",
			expectedResult: 2,
		},
		{
			name:           "priority",
			expression:     "(2+2)*2",
			expectedResult: 8,
		},
		{
			name:           "priority",
			expression:     "2+2*2",
			expectedResult: 6,
		},
		{
			name:           "/",
			expression:     "1/2",
			expectedResult: 0.5,
		},
		{
			name:           "justNumber",
			expression:     "3",
			expectedResult: 3,
		},
		{
			name:           "longExpression",
			expression:     "(4+3+2)/(1+2) * 1 / 3",
			expectedResult: 1,
		},
		{
			name:           "longExpression2",
			expression:     "((1+4) * (1+2) +1) *4",
			expectedResult: 64,
		},
		{
			name:           "expressionWithSpaces",
			expression:     "1+1 * 2",
			expectedResult: 3,
		},
	}

	for _, testCase := range testCasesSuccess {
		t.Run(testCase.name, func(t *testing.T) {
			val, err := calculation.Calc(testCase.expression)
			if err != nil {
				t.Fatalf("successful case %s returns error", testCase.expression)
			}
			if val != testCase.expectedResult {
				t.Fatalf("%f should be equal %f", val, testCase.expectedResult)
			}
		})
	}

	testCasesFail := []struct {
		name        string
		expression  string
		expectedErr error
	}{
		{
			name:       "simple",
			expression: "1+1*",
		},
		{
			name:       "priority",
			expression: "2+2**2",
		},
		{
			name:       "priority",
			expression: "((2+2-*(2",
		},
		{
			name:       "/",
			expression: "",
		},
		{
			name:       "unknownOperation",
			expression: "1**2",
		},
		{
			name:       "unknownOperation2",
			expression: "6^2",
		},
		{
			name:       "justOperation",
			expression: "-",
		},
		{
			name:       "zeroDivision",
			expression: "1/0",
		},
	}

	for _, testCase := range testCasesFail {
		t.Run(testCase.name, func(t *testing.T) {
			val, err := calculation.Calc(testCase.expression)
			if err == nil {
				t.Fatalf("expression %s is invalid but result  %f was obtained", testCase.expression, val)
			}
		})
	}
}
