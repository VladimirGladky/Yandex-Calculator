package orchestratorapp

import (
	"github.com/VladimirGladky/FinalTaskFirstSprint/internal/orchestrator/transport/http"
)

type App struct {
	orch *http.Orchestrator
}

func New(orchestrator *http.Orchestrator) *App {
	return &App{
		orch: orchestrator,
	}
}

func (a *App) Run() error {
	if err := a.orch.Run(); err != nil {
		return err
	}
	return nil
}
