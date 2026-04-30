package service

import (
	"context"
	"fmt"

	orderdto "github.com/lootarola/ai-incident-response-challenge/internal/order/handler/rest/dto"
	telemetryotel "github.com/lootarola/ai-incident-response-challenge/pkg/telemetry/otel"
	"github.com/lootarola/ai-incident-response-challenge/pkg/types"
)

var validTransitions = map[types.OrderStatus][]types.OrderStatus{
	types.StatusPending:   {types.StatusConfirmed, types.StatusCancelled},
	types.StatusConfirmed: {types.StatusShipped, types.StatusCancelled},
	types.StatusShipped:   {types.StatusDelivered},
}

func (s *Service) Update(ctx context.Context, id string, req orderdto.UpdateRequest) (*types.Order, error) {
	ctx = telemetryotel.WithDomain(ctx, "order")
	ctx = telemetryotel.WithOperation(ctx, "order.update")
	ctx, span := s.tracer.Start(ctx, "order.update")
	defer span.End()

	existing, err := s.repo.Get(ctx, id)
	if err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("get order for update: %w", mapRepoErr(err))
	}

	allowed := validTransitions[existing.Status]
	valid := false
	for _, s := range allowed {
		if s == req.Status {
			valid = true
			break
		}
	}
	if !valid {
		return nil, fmt.Errorf("transition %s→%s: %w", existing.Status, req.Status, types.ErrInvalidStatus)
	}

	existing.Status = req.Status
	if err := s.repo.Update(ctx, existing); err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("update order: %w", mapRepoErr(err))
	}

	s.logger.Info(ctx, "order updated",
		telemetryotel.KV("domain", telemetryotel.Domain(ctx)),
		telemetryotel.KV("operation", telemetryotel.Operation(ctx)),
		telemetryotel.KV("order_id", id),
		telemetryotel.KV("status", string(req.Status)),
	)
	return existing, nil
}
