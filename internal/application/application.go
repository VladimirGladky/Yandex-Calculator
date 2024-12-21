package application

import (
	"FinalTaskFirstSprint/pkg/calculation"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"strconv"
)

type Config struct {
	Addr string
}

func ConfigFromEnv() *Config {
	config := new(Config)
	config.Addr = os.Getenv("PORT")
	if config.Addr == "" {
		config.Addr = "4040"
	}
	return config
}

type Application struct {
	config *Config
}

func New() *Application {
	return &Application{
		config: ConfigFromEnv(),
	}
}

type Request struct {
	Expression string `json:"expression"`
}

type BadResponse struct {
	Error string `json:"error"`
}

type GoodResponse struct {
	Result string `json:"result"`
}

func CalcHandler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if rec := recover(); rec != nil {
			w.WriteHeader(http.StatusInternalServerError)
			err0 := json.NewEncoder(w).Encode(BadResponse{Error: "Internal server error1"})
			if err0 != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"error": "Internal server error2"}`))
			}
			return
		}
	}()

	if r.Method != http.MethodPost {
		//405
		w.WriteHeader(http.StatusMethodNotAllowed)
		err10 := json.NewEncoder(w).Encode(BadResponse{Error: "You can use only POST method"})
		if err10 != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error": "Internal server error3"}`))
		}
		return
	}

	request := new(Request)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			err1 := json.NewEncoder(w).Encode(BadResponse{Error: "Internal server error4"})
			if err1 != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"error": "Internal server error5"}`))
			}
			return
		}
	}(r.Body)
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		//400
		w.WriteHeader(http.StatusBadRequest)
		err2 := json.NewEncoder(w).Encode(BadResponse{Error: "Bad request"})
		if err2 != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error": "Internal server error6"}`))
		}
		return
	}

	result, err := calculation.Calc(request.Expression)
	if err != nil {
		if errors.Is(err, calculation.ErrInvalidExpression) {
			//422
			w.WriteHeader(http.StatusUnprocessableEntity)
			err3 := json.NewEncoder(w).Encode(BadResponse{Error: "Expression is not valid"})
			if err3 != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"error": "Internal server error7"}`))
			}
		} else {
			w.WriteHeader(http.StatusUnprocessableEntity)
			err4 := json.NewEncoder(w).Encode(BadResponse{Error: "You can't divide by zero"})
			if err4 != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"error": "Internal server error8"}`))
			}
		}
	} else {
		err5 := json.NewEncoder(w).Encode(GoodResponse{Result: strconv.FormatFloat(result, 'f', 8, 64)})
		if err5 != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error": "Internal server error9"}`))
		}
	}
}

func (a *Application) RunServer() error {
	http.HandleFunc("/api/v1/calculate", CalcHandler)
	return http.ListenAndServe(":"+a.config.Addr, nil)
}
