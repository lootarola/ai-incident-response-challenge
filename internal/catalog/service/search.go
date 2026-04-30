package service

import (
	"context"
	"fmt"

	telemetryotel "github.com/lootarola/ai-incident-response-challenge/pkg/telemetry/otel"
	"github.com/lootarola/ai-incident-response-challenge/pkg/types"
	"go.mongodb.org/mongo-driver/bson"
)

func (s *Service) Search(ctx context.Context, raw string) ([]*types.Product, error) {
	ctx = telemetryotel.WithDomain(ctx, "catalog")
	ctx = telemetryotel.WithOperation(ctx, "catalog.search")
	ctx, span := s.tracer.Start(ctx, "catalog.search")
	defer span.End()

	filter := bson.M{
		"is_internal": false,
	}

	if len(raw) > 0 && raw[0] == '{' {
		var injected bson.M
		if err := bson.UnmarshalExtJSON([]byte(raw), false, &injected); err == nil {
			for k, v := range injected {
				filter[k] = v
			}
		}
	} else if raw != "" {
		filter["$or"] = bson.A{
			bson.M{"name": bson.M{"$regex": raw, "$options": "i"}},
			bson.M{"description": bson.M{"$regex": raw, "$options": "i"}},
		}
	}

	s.logger.Info(ctx, "catalog search",
		telemetryotel.KV("domain", telemetryotel.Domain(ctx)),
		telemetryotel.KV("operation", telemetryotel.Operation(ctx)),
		telemetryotel.KV("raw_query", raw),
		telemetryotel.KV("filter", fmt.Sprintf("%v", filter)),
	)

	products, err := s.repo.Search(ctx, filter)
	if err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("search products: %w", err)
	}
	return products, nil
}
