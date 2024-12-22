package application_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/VladimirGladky/FinalTaskFirstSprint/internal/application"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCalc(t *testing.T) {
	testCasesSuccess := []struct {
		name           string
		method         string
		requestBody    application.Request
		expectedCode   int
		expectedResult application.GoodResponse
	}{
		{
			name:           "simple",
			method:         "POST",
			requestBody:    application.Request{Expression: "1+1"},
			expectedCode:   200,
			expectedResult: application.GoodResponse{Result: "2.00000000"},
		},
		{
			name:           "priority",
			method:         "POST",
			requestBody:    application.Request{Expression: "1+2*2"},
			expectedCode:   200,
			expectedResult: application.GoodResponse{Result: "5.00000000"},
		},
		{
			name:           "priority2",
			method:         "POST",
			requestBody:    application.Request{Expression: "(1+1)*2"},
			expectedCode:   200,
			expectedResult: application.GoodResponse{Result: "4.00000000"},
		},
	}

	for _, testCase := range testCasesSuccess {
		t.Run(testCase.name, func(t *testing.T) {
			data := testCase.requestBody
			jsonData, err := json.Marshal(data)
			if err != nil {
				fmt.Println("Error marshaling JSON:", err)
				return
			}
			req := httptest.NewRequest(http.MethodPost, "/api/v1/calculate", bytes.NewBuffer(jsonData))
			w := httptest.NewRecorder()
			application.CalcHandler(w, req)
			res := w.Result()
			defer res.Body.Close()
			body, err := io.ReadAll(res.Body)
			if err != nil {
				t.Fatalf("Error reading response body: %v", err)
			}
			var response application.GoodResponse
			if err1 := json.Unmarshal(body, &response); err1 != nil {
				t.Fatalf("Error unmarshaling JSON: %v", err)
			}
			if response.Result != testCase.expectedResult.Result {
				t.Errorf("Expected %v but got %v", testCase.expectedResult.Result, response.Result)
			}
			if res.StatusCode != testCase.expectedCode {
				t.Errorf("wrong status code")
			}
			if req.Method != testCase.method {
				t.Errorf("wrong method")
			}
		})
	}

	testCasesFail := []struct {
		name           string
		method         string
		requestBody    application.Request
		expectedCode   int
		expectedResult application.BadResponse
	}{
		{
			name:           "manyOperations",
			method:         "POST",
			requestBody:    application.Request{Expression: "1+1+"},
			expectedCode:   422,
			expectedResult: application.BadResponse{Error: "Expression is not valid"},
		},
		{
			name:           "extraBracket",
			method:         "POST",
			requestBody:    application.Request{Expression: "2*(10+9"},
			expectedCode:   422,
			expectedResult: application.BadResponse{Error: "Expression is not valid"},
		},
		{
			name:           "not numbs",
			method:         "POST",
			requestBody:    application.Request{Expression: "not numbs"},
			expectedCode:   422,
			expectedResult: application.BadResponse{Error: "Expression is not valid"},
		},
	}

	for _, testCase := range testCasesFail {
		t.Run(testCase.name, func(t *testing.T) {
			data := testCase.requestBody
			jsonData, err := json.Marshal(data)
			if err != nil {
				fmt.Println("Error marshaling JSON:", err)
				return
			}
			req := httptest.NewRequest(http.MethodPost, "/api/v1/calculate", bytes.NewBuffer(jsonData))
			w := httptest.NewRecorder()
			application.CalcHandler(w, req)
			res := w.Result()
			defer res.Body.Close()
			body, err := io.ReadAll(res.Body)
			if err != nil {
				t.Fatalf("Error reading response body: %v", err)
			}
			var response application.BadResponse
			if err1 := json.Unmarshal(body, &response); err1 != nil {
				t.Fatalf("Error unmarshaling JSON: %v", err)
			}
			if response.Error != testCase.expectedResult.Error {
				t.Errorf("Expected %v but got %v", testCase.expectedResult.Error, response.Error)
			}
			if res.StatusCode != testCase.expectedCode {
				t.Errorf("wrong status code")
			}
			if req.Method != testCase.method {
				t.Errorf("wrong method")
			}
		})
	}
}

// тест для ошибки 405
func TestCalcHandlerWrongMethodCase(t *testing.T) {
	jsonData, err := json.Marshal(application.Request{Expression: "1+1 * 2"})
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}
	req := httptest.NewRequest(http.MethodGet, "/api/v1/calculate", bytes.NewBuffer(jsonData))
	w := httptest.NewRecorder()
	application.CalcHandler(w, req)
	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("wrong status code")
	}
}

// тест для ошибки 400
func TestCalcHandlerBadRequest(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/v1/calculate", bytes.NewBuffer([]byte("fjdng")))
	w := httptest.NewRecorder()
	application.CalcHandler(w, req)
	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("wrong status code")
	}
}
