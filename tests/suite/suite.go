package suite

import (
	"github.com/VladimirGladky/FinalTaskFirstSprint/internal/orchestrator/transport/http"
	"testing"
)

type Suite struct {
	*testing.T
	Orchestrator *http.Orchestrator
}

func New(t *testing.T) *Suite {
	//t.Helper()
	//t.Parallel()
	//ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	//ctx, err := logger.New(ctx)
	//require.NoError(t, err)
	//cfg := config.Config{
	//	OrchestratorPort:      "0",
	//	OrchestratorHost:      "localhost",
	//	GrpcPort:              "123",
	//	GrpcHost:              "localhost",
	//	ComputingPower:        3,
	//	TimeAdditionMs:        200,
	//	TimeSubtractionMs:     200,
	//	TimeMultiplicationsMs: 200,
	//	TimeDivisionsMs:       200,
	//}
	//require.NoError(t, err)
	//cfg.GrpcHost = "localhost"
	//cfg.GrpcPort = "0"
	//cfg.OrchestratorHost = "localhost"
	//cfg.OrchestratorPort = "123"
	//o := http.New(ctx, &mock.Service{}, &cfg)
	//t.Cleanup(func() {
	//	t.Helper()
	//	cancel()
	//})
	//return &Suite{
	//	T:            t,
	//	Orchestrator: o,
	//}
	return nil
}
