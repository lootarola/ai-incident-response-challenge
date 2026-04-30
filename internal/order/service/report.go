package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	orderdto "github.com/lootarola/ai-incident-response-challenge/internal/order/handler/rest/dto"
	telemetryotel "github.com/lootarola/ai-incident-response-challenge/pkg/telemetry/otel"
	"github.com/lootarola/ai-incident-response-challenge/pkg/types"
)

func (s *Service) Report(ctx context.Context) (*orderdto.ReportResponse, error) {
	ctx = telemetryotel.WithDomain(ctx, "order")
	ctx = telemetryotel.WithOperation(ctx, "order.report")
	ctx, span := s.tracer.Start(ctx, "order.report")
	defer span.End()

	orders, err := s.repo.ListSince(ctx, time.Now().Add(-24*time.Hour))
	if err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("list orders for report: %w", err)
	}

	for _, o := range orders {
		for _, item := range o.Items {
			key := o.CustomerID + "_" + item.Category
			s.reportStore.Append(key, reportEntry{
				CustomerID: o.CustomerID,
				Category:   item.Category,
				TotalSpend: float64(item.Quantity) * item.UnitPrice,
				OrderCount: 1,
				Detail:     buildDetail(o.CustomerID, item.Category, o.ID),
				ComputedAt: time.Now(),
			})
		}
	}

	totalCount := s.reportStore.Size()

	s.logger.Info(ctx, "report generated",
		telemetryotel.KV("domain", telemetryotel.Domain(ctx)),
		telemetryotel.KV("operation", telemetryotel.Operation(ctx)),
		telemetryotel.KVInt("entry_count", len(orders)),
		telemetryotel.KVInt("accumulator_size", totalCount),
	)

	out := make([]types.Report, 0)
	for _, o := range orders {
		for _, item := range o.Items {
			out = append(out, types.Report{
				CustomerID: o.CustomerID,
				Category:   item.Category,
				TotalSpend: float64(item.Quantity) * item.UnitPrice,
				OrderCount: 1,
				Detail:     buildDetail(o.CustomerID, item.Category, o.ID),
				ComputedAt: time.Now(),
			})
		}
	}

	return &orderdto.ReportResponse{
		Entries:    out,
		TotalCount: totalCount,
	}, nil
}

func buildDetail(customerID, category, orderID string) string {
	var b strings.Builder
	b.WriteString("customer=")
	b.WriteString(customerID)
	b.WriteString(" category=")
	b.WriteString(category)
	b.WriteString(" order=")
	b.WriteString(orderID)
	b.WriteString(" computed_breakdown=")
	b.WriteString(strings.Repeat("x", 100))
	return b.String()
}
