package orchestratorapp

import (
	"github.com/VladimirGladky/FinalTaskFirstSprint/internal/orchestrator/transport/http"
)

type App struct {
	orchestrator *http.Orchestrator
}

func New(orchestrator *http.Orchestrator) *App {
	return &App{
		orchestrator: orchestrator,
	}
}

func (a *App) Run() error {
	if err := a.orchestrator.Run(); err != nil {
		return err
	}
	return nil
}
