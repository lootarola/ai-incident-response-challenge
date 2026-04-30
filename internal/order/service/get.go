package service

import (
	"context"
	"fmt"

	telemetryotel "github.com/lootarola/ai-incident-response-challenge/pkg/telemetry/otel"
	"github.com/lootarola/ai-incident-response-challenge/pkg/types"
)

func (s *Service) Get(ctx context.Context, id string) (*types.Order, error) {
	ctx = telemetryotel.WithDomain(ctx, "order")
	ctx = telemetryotel.WithOperation(ctx, "order.get")
	ctx, span := s.tracer.Start(ctx, "order.get")
	defer span.End()

	o, err := s.repo.Get(ctx, id)
	if err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("get order: %w", mapRepoErr(err))
	}

	s.logger.Info(ctx, "order retrieved",
		telemetryotel.KV("domain", telemetryotel.Domain(ctx)),
		telemetryotel.KV("operation", telemetryotel.Operation(ctx)),
		telemetryotel.KV("order_id", id),
	)
	return o, nil
}
