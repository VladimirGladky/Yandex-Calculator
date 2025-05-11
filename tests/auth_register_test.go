package tests

import (
	_ "net/http/httptest"
)

//func TestAuthRegisterGood(t *testing.T) {
//	st := suite.New(t)
//
//	cases := []struct {
//		name         string
//		registerData models.RegisterRequest
//		method       string
//		statusCode   int
//		expectedRes  string
//	}{
//		{
//			name: "test1",
//			registerData: models.RegisterRequest{
//				Username: "test1",
//				Password: "943828393!",
//			},
//			method:      "POST",
//			statusCode:  http.StatusOK,
//			expectedRes: "success",
//		},
//		{
//			name: "test2",
//			registerData: models.RegisterRequest{
//				Username: "test2",
//				Password: "qwerty1234#",
//			},
//			method:      "POST",
//			statusCode:  http.StatusUnprocessableEntity,
//			expectedRes: "success",
//		},
//		{
//			name: "test3",
//			registerData: models.RegisterRequest{
//				Username: "test3",
//				Password: "groovy1234%",
//			},
//			method:      "POST",
//			statusCode:  http.StatusMethodNotAllowed,
//			expectedRes: "success",
//		},
//	}
//
//	for _, tc := range cases {
//		t.Run(tc.name, func(t *testing.T) {
//			jsonData, err := json.Marshal(tc.registerData)
//			require.NoError(t, err)
//
//			req := httptest.NewRequest(tc.method, "http://localhost:123/api/v1/register", bytes.NewBuffer(jsonData))
//
//			w := httptest.NewRecorder()
//			http2.RegisterHandler(st.Orchestrator)(w, req)
//
//			resp := w.Result()
//			defer resp.Body.Close()
//
//			if resp.StatusCode != tc.statusCode {
//				t.Errorf("Expected status code %d, got %d", tc.statusCode, resp.StatusCode)
//			}
//
//			var response models.BadResponse
//			err = json.NewDecoder(resp.Body).Decode(&response)
//			t.Errorf("Expected response %s, got %s", tc.expectedRes, response.Error)
//			require.NoError(t, err)
//		})
//	}
//}
//
//func TestAuthRegisterBad(t *testing.T) {
//	st := suite.New(t)
//	err := st.Orchestrator.Run()
//	require.NoError(t, err)
//
//	cases := []struct {
//		name         string
//		registerData models.RegisterRequest
//		method       string
//		statusCode   int
//		expectedRes  string
//	}{
//		{
//			name: "test1",
//			registerData: models.RegisterRequest{
//				Username: "test1",
//				Password: "943828393!",
//			},
//			method:      "POST",
//			statusCode:  http.StatusOK,
//			expectedRes: "success",
//		},
//	}
//
//	for _, tc := range cases {
//		t.Run(tc.name, func(t *testing.T) {
//			jsonData, err := json.Marshal(tc.registerData)
//			require.NoError(t, err)
//
//			req := httptest.NewRequest(tc.method, "http://localhost:123/api/v1/register", bytes.NewBuffer(jsonData))
//			req.Header.Set("Content-Type", "application/json")
//
//			resp, err := http.DefaultClient.Do(req)
//			require.NoError(t, err)
//
//			if resp.StatusCode != tc.statusCode {
//				t.Errorf("Expected status code %d, got %d", tc.statusCode, resp.StatusCode)
//			}
//
//			var response models.RegisterGoodResponse
//			err = json.NewDecoder(resp.Body).Decode(&response)
//			require.NoError(t, err)
//			require.Equal(t, tc.expectedRes, response.Status)
//		})
//	}
//}
