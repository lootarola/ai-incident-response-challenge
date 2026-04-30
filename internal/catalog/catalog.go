package catalog

import (
	"fmt"

	catalogrest "github.com/lootarola/ai-incident-response-challenge/internal/catalog/handler/rest"
	catalogservice "github.com/lootarola/ai-incident-response-challenge/internal/catalog/service"
	"github.com/lootarola/ai-incident-response-challenge/pkg/database"
	"github.com/lootarola/ai-incident-response-challenge/pkg/server"
	telemetryotel "github.com/lootarola/ai-incident-response-challenge/pkg/telemetry/otel"
)

type Repository = database.CatalogRepository

type Module struct {
	Handler server.Handler
}

func New(repo Repository, scope *telemetryotel.Scope) (*Module, error) {
	svc, err := catalogservice.NewService(repo, scope)
	if err != nil {
		return nil, fmt.Errorf("catalog service: %w", err)
	}
	h, err := catalogrest.NewHandler(svc)
	if err != nil {
		return nil, fmt.Errorf("catalog handler: %w", err)
	}
	return &Module{Handler: h}, nil
}
