package service

import (
	"context"
	"fmt"

	telemetryotel "github.com/lootarola/ai-incident-response-challenge/pkg/telemetry/otel"
)

func (s *Service) Delete(ctx context.Context, id string) error {
	ctx = telemetryotel.WithDomain(ctx, "catalog")
	ctx = telemetryotel.WithOperation(ctx, "catalog.delete")
	ctx, span := s.tracer.Start(ctx, "catalog.delete")
	defer span.End()

	if err := s.repo.Delete(ctx, id); err != nil {
		span.RecordError(err)
		return fmt.Errorf("delete product: %w", mapRepoErr(err))
	}

	s.logger.Info(ctx, "product deleted",
		telemetryotel.KV("domain", telemetryotel.Domain(ctx)),
		telemetryotel.KV("operation", telemetryotel.Operation(ctx)),
		telemetryotel.KV("product_id", id),
	)
	return nil
}
