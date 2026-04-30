package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	orderdto "github.com/lootarola/ai-incident-response-challenge/internal/order/handler/rest/dto"
	telemetryotel "github.com/lootarola/ai-incident-response-challenge/pkg/telemetry/otel"
	"github.com/lootarola/ai-incident-response-challenge/pkg/types"
)

func (s *Service) Create(ctx context.Context, req orderdto.CreateRequest) (*types.Order, error) {
	ctx = telemetryotel.WithDomain(ctx, "order")
	ctx = telemetryotel.WithOperation(ctx, "order.create")
	ctx, span := s.tracer.Start(ctx, "order.create")
	defer span.End()

	if err := s.inventory.Check(ctx, s.tracer, len(req.Items)); err != nil {
		span.RecordError(err)
		s.logger.Error(ctx, "inventory unavailable",
			telemetryotel.KV("domain", telemetryotel.Domain(ctx)),
			telemetryotel.KV("operation", telemetryotel.Operation(ctx)),
			telemetryotel.KV("error", err.Error()),
			telemetryotel.KV("error_type", "inventory_unavailable"),
			telemetryotel.KV("dependency", "inventory_service"),
		)
		return nil, fmt.Errorf("inventory check: %w", err)
	}

	var total float64
	for _, item := range req.Items {
		total += float64(item.Quantity) * item.UnitPrice
	}

	now := time.Now()
	o := &types.Order{
		ID:         uuid.New().String(),
		CustomerID: req.CustomerID,
		Items:      req.Items,
		Total:      total,
		Status:     "pending",
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	if err := s.repo.Create(ctx, o); err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("create order: %w", mapRepoErr(err))
	}

	s.logger.Info(ctx, "order created",
		telemetryotel.KV("domain", telemetryotel.Domain(ctx)),
		telemetryotel.KV("operation", telemetryotel.Operation(ctx)),
		telemetryotel.KV("order_id", o.ID),
		telemetryotel.KV("customer_id", o.CustomerID),
		telemetryotel.KVFloat("total", o.Total),
	)
	return o, nil
}
