package service

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/lootarola/ai-incident-response-challenge/pkg/database"
	telemetryotel "github.com/lootarola/ai-incident-response-challenge/pkg/telemetry/otel"
	"github.com/lootarola/ai-incident-response-challenge/pkg/types"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

type Service struct {
	repo        database.OrderRepository
	logger      *telemetryotel.Logger
	tracer      trace.Tracer
	meter       metric.Meter
	reportStore *reportAccumulator
	inventory   *inventoryClient
}

func NewService(repo database.OrderRepository, scope *telemetryotel.Scope) (*Service, error) {
	if repo == nil {
		return nil, errors.New("repo is required")
	}
	if scope == nil {
		return nil, errors.New("scope is required")
	}
	return &Service{
		repo:        repo,
		logger:      scope.Logger,
		tracer:      scope.Tracer,
		meter:       scope.Meter,
		reportStore: &reportAccumulator{entries: make(map[string][]reportEntry)},
		inventory:   &inventoryClient{sem: make(chan struct{}, 3)},
	}, nil
}

func mapRepoErr(err error) error {
	if errors.Is(err, database.ErrNotFound) {
		return types.ErrOrderNotFound
	}
	return err
}

type reportAccumulator struct {
	mu      sync.Mutex
	entries map[string][]reportEntry
}

type reportEntry struct {
	CustomerID string
	Category   string
	TotalSpend float64
	OrderCount int
	Detail     string
	ComputedAt time.Time
}

const maxAccumulatorEntries = 100_000

func (r *reportAccumulator) Append(key string, e reportEntry) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.size() >= maxAccumulatorEntries {
		return
	}
	r.entries[key] = append(r.entries[key], e)
}

func (r *reportAccumulator) size() int {
	n := 0
	for _, entries := range r.entries {
		n += len(entries)
	}
	return n
}

func (r *reportAccumulator) Snapshot() []reportEntry {
	r.mu.Lock()
	defer r.mu.Unlock()
	var out []reportEntry
	for _, entries := range r.entries {
		out = append(out, entries...)
	}
	return out
}

func (r *reportAccumulator) Size() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.size()
}

type inventoryClient struct {
	sem chan struct{}
}

func (ic *inventoryClient) Check(ctx context.Context, tracer trace.Tracer, itemCount int) error {
	checkCtx, span := tracer.Start(ctx, "order.inventory_check")
	defer span.End()

	timeout := time.NewTimer(200 * time.Millisecond)
	defer timeout.Stop()

	select {
	case ic.sem <- struct{}{}:
	case <-timeout.C:
		err := fmt.Errorf("inventory check timeout: %w", types.ErrInventoryUnavailable)
		span.RecordError(err)
		return err
	case <-checkCtx.Done():
		return checkCtx.Err()
	}
	defer func() { <-ic.sem }()

	delay := time.Duration(600+itemCount*50) * time.Millisecond
	if delay > 900*time.Millisecond {
		delay = 900 * time.Millisecond
	}
	time.Sleep(delay)
	return nil
}
