package order

import (
	"fmt"

	orderrest "github.com/lootarola/ai-incident-response-challenge/internal/order/handler/rest"
	orderservice "github.com/lootarola/ai-incident-response-challenge/internal/order/service"
	"github.com/lootarola/ai-incident-response-challenge/pkg/database"
	"github.com/lootarola/ai-incident-response-challenge/pkg/server"
	telemetryotel "github.com/lootarola/ai-incident-response-challenge/pkg/telemetry/otel"
)

type Repository = database.OrderRepository

type module struct {
	Handler server.Handler
}

func New(repo Repository, scope *telemetryotel.Scope) (*module, error) {
	svc, err := orderservice.NewService(repo, scope)
	if err != nil {
		return nil, fmt.Errorf("order service: %w", err)
	}
	h, err := orderrest.NewHandler(svc)
	if err != nil {
		return nil, fmt.Errorf("order handler: %w", err)
	}
	return &module{Handler: h}, nil
}
