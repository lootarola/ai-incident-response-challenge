package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	catalogdto "github.com/lootarola/ai-incident-response-challenge/internal/catalog/handler/rest/dto"
	telemetryotel "github.com/lootarola/ai-incident-response-challenge/pkg/telemetry/otel"
	"github.com/lootarola/ai-incident-response-challenge/pkg/types"
)

func (s *Service) Create(ctx context.Context, req catalogdto.CreateRequest) (*types.Product, error) {
	ctx = telemetryotel.WithDomain(ctx, "catalog")
	ctx = telemetryotel.WithOperation(ctx, "catalog.create")
	ctx, span := s.tracer.Start(ctx, "catalog.create")
	defer span.End()

	now := time.Now()
	p := &types.Product{
		ID:          uuid.New().String(),
		Name:        req.Name,
		Description: req.Description,
		Category:    req.Category,
		Price:       req.Price,
		Inventory:   req.Inventory,
		IsInternal:  req.IsInternal,
		Specs:       req.Specs,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := s.repo.Create(ctx, p); err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("create product: %w", mapRepoErr(err))
	}

	s.logger.Info(ctx, "product created",
		telemetryotel.KV("domain", telemetryotel.Domain(ctx)),
		telemetryotel.KV("operation", telemetryotel.Operation(ctx)),
		telemetryotel.KV("product_id", p.ID),
		telemetryotel.KV("category", p.Category),
	)
	return p, nil
}
