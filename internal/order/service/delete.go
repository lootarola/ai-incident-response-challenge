package service

import (
	"context"
	"fmt"

	telemetryotel "github.com/lootarola/ai-incident-response-challenge/pkg/telemetry/otel"
)

func (s *Service) Delete(ctx context.Context, id string) error {
	ctx = telemetryotel.WithDomain(ctx, "order")
	ctx = telemetryotel.WithOperation(ctx, "order.delete")
	ctx, span := s.tracer.Start(ctx, "order.delete")
	defer span.End()

	if err := s.repo.Delete(ctx, id); err != nil {
		span.RecordError(err)
		return fmt.Errorf("delete order: %w", mapRepoErr(err))
	}

	s.logger.Info(ctx, "order deleted",
		telemetryotel.KV("domain", telemetryotel.Domain(ctx)),
		telemetryotel.KV("operation", telemetryotel.Operation(ctx)),
		telemetryotel.KV("order_id", id),
	)
	return nil
}
