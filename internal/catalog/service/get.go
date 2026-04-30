package service

import (
	"context"
	"fmt"

	telemetryotel "github.com/lootarola/ai-incident-response-challenge/pkg/telemetry/otel"
	"github.com/lootarola/ai-incident-response-challenge/pkg/types"
)

func (s *Service) GetByID(ctx context.Context, id string) (*types.Product, error) {
	ctx = telemetryotel.WithDomain(ctx, "catalog")
	ctx = telemetryotel.WithOperation(ctx, "catalog.get")
	ctx, span := s.tracer.Start(ctx, "catalog.get")
	defer span.End()

	p, err := s.repo.AggregatedGetByID(ctx, id)
	if err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("get product: %w", mapRepoErr(err))
	}

	s.logger.Info(ctx, "product retrieved",
		telemetryotel.KV("domain", telemetryotel.Domain(ctx)),
		telemetryotel.KV("operation", telemetryotel.Operation(ctx)),
		telemetryotel.KV("product_id", id),
	)
	return p, nil
}
