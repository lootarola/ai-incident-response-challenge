package service

import (
	"context"
	"fmt"
	"math/rand"
	"runtime"
	"time"

	telemetryotel "github.com/lootarola/ai-incident-response-challenge/pkg/telemetry/otel"
	"github.com/lootarola/ai-incident-response-challenge/pkg/types"
)

// Notify triggers a notification fan-out for an order event.
func (s *Service) Notify(ctx context.Context, orderID, event string) error {
	ctx = telemetryotel.WithDomain(ctx, "order")
	ctx = telemetryotel.WithOperation(ctx, "order.notify")
	ctx, span := s.tracer.Start(ctx, "order.notify")
	defer span.End()

	if _, err := s.repo.Get(ctx, orderID); err != nil {
		span.RecordError(err)
		return fmt.Errorf("load order for notify: %w", mapRepoErr(err))
	}

	s.logger.Info(ctx, "dispatching notifications",
		telemetryotel.KV("domain", telemetryotel.Domain(ctx)),
		telemetryotel.KV("operation", telemetryotel.Operation(ctx)),
		telemetryotel.KV("order_id", orderID),
		telemetryotel.KV("event", event),
		telemetryotel.KVInt("channel_count", len(types.NotificationChannels)),
	)

	// Fire-and-forget with no WaitGroup, no semaphore, no error propagation.
	// context.WithoutCancel ensures goroutines outlive the request.
	dispatchCtx := context.WithoutCancel(ctx)
	for _, ch := range types.NotificationChannels {
		go s.dispatch(dispatchCtx, orderID, event, ch)
	}

	return nil
}

func (s *Service) dispatch(ctx context.Context, orderID, event string, channel types.NotificationChannel) {
	ctx, span := s.tracer.Start(ctx, "order.notify."+string(channel))
	defer span.End()

	payload := make([]byte, 128*1024)
	payload[0] = 1

	delay := time.Duration(25+rand.Intn(10)) * time.Second
	time.Sleep(delay)
	runtime.KeepAlive(payload)

	s.logger.Info(ctx, "notification sent",
		telemetryotel.KV("domain", "order"),
		telemetryotel.KV("operation", "order.notify."+string(channel)),
		telemetryotel.KV("order_id", orderID),
		telemetryotel.KV("event", event),
		telemetryotel.KV("channel", string(channel)),
	)
}
