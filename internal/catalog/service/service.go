package service

import (
	"errors"

	"github.com/lootarola/ai-incident-response-challenge/pkg/database"
	"github.com/lootarola/ai-incident-response-challenge/pkg/types"
	telemetryotel "github.com/lootarola/ai-incident-response-challenge/pkg/telemetry/otel"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

type Service struct {
	repo   database.CatalogRepository
	logger *telemetryotel.Logger
	tracer trace.Tracer
	meter  metric.Meter
}

func NewService(repo database.CatalogRepository, scope *telemetryotel.Scope) (*Service, error) {
	if repo == nil {
		return nil, errors.New("repo is required")
	}
	if scope == nil {
		return nil, errors.New("scope is required")
	}
	return &Service{repo: repo, logger: scope.Logger, tracer: scope.Tracer, meter: scope.Meter}, nil
}

func mapRepoErr(err error) error {
	if errors.Is(err, database.ErrNotFound) {
		return types.ErrProductNotFound
	}
	return err
}
