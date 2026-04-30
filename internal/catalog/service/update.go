package service

import (
	"context"
	"fmt"

	catalogdto "github.com/lootarola/ai-incident-response-challenge/internal/catalog/handler/rest/dto"
	telemetryotel "github.com/lootarola/ai-incident-response-challenge/pkg/telemetry/otel"
	"github.com/lootarola/ai-incident-response-challenge/pkg/types"
)

func (s *Service) Update(ctx context.Context, id string, req catalogdto.UpdateRequest) (*types.Product, error) {
	ctx = telemetryotel.WithDomain(ctx, "catalog")
	ctx = telemetryotel.WithOperation(ctx, "catalog.update")
	ctx, span := s.tracer.Start(ctx, "catalog.update")
	defer span.End()

	existing, err := s.repo.AggregatedGetByID(ctx, id)
	if err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("get product for update: %w", mapRepoErr(err))
	}

	if req.Name != "" {
		existing.Name = req.Name
	}
	if req.Description != "" {
		existing.Description = req.Description
	}
	if req.Category != "" {
		existing.Category = req.Category
	}
	if req.Price > 0 {
		existing.Price = req.Price
	}
	if req.Inventory > 0 {
		existing.Inventory = req.Inventory
	}
	if req.Specs != nil {
		existing.Specs = req.Specs
	}

	if err := s.repo.Update(ctx, existing); err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("update product: %w", mapRepoErr(err))
	}

	s.logger.Info(ctx, "product updated",
		telemetryotel.KV("domain", telemetryotel.Domain(ctx)),
		telemetryotel.KV("operation", telemetryotel.Operation(ctx)),
		telemetryotel.KV("product_id", id),
	)
	return existing, nil
}
