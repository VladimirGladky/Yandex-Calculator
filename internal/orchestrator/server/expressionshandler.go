package server

import (
	"encoding/json"
	models2 "github.com/VladimirGladky/FinalTaskFirstSprint/internal/models"
	"github.com/VladimirGladky/FinalTaskFirstSprint/pkg/logger"
	"net/http"
)

func ExpressionsHandler(o *Orchestrator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				logger.GetLoggerFromCtx(o.Ctx).Error(o.Ctx, "Internal server error")
				w.WriteHeader(http.StatusInternalServerError)
				err0 := json.NewEncoder(w).Encode(models2.BadResponse{Error: "Internal server error1"})
				if err0 != nil {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte(`{"error": "Internal server error2"}`))
				}
				return
			}
		}()
		if r.Method != http.MethodGet {
			//405
			logger.GetLoggerFromCtx(o.Ctx).Info(o.Ctx, "You can use only GET method")
			w.WriteHeader(http.StatusMethodNotAllowed)
			err10 := json.NewEncoder(w).Encode(models2.BadResponse{Error: "You can use only GET method"})
			if err10 != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"error": "Internal server error3"}`))
			}
			return
		}
		expressions := make([]models2.Expression, 0, len(o.ExpressionsMap))
		for _, v := range o.ExpressionsMap {
			if v.Ast != nil && v.Ast.IsLeaf {
				v.Status = "done"
				v.Result = v.Ast.Value
			}
			expressions = append(expressions, *v)
		}
		err := json.NewEncoder(w).Encode(models2.Expressions{Expressions: expressions})
		if err != nil {
			logger.GetLoggerFromCtx(o.Ctx).Error(o.Ctx, "Internal server error")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error": "Internal server error4"}`))
			return
		}
		w.WriteHeader(http.StatusOK)
		logger.GetLoggerFromCtx(o.Ctx).Info(o.Ctx, "Expressions sent")
	}
}
