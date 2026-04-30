package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	orderdto "github.com/lootarola/ai-incident-response-challenge/internal/order/handler/rest/dto"
	"github.com/lootarola/ai-incident-response-challenge/internal/order/service"
	"github.com/lootarola/ai-incident-response-challenge/pkg/database/mock"
	telemetryotel "github.com/lootarola/ai-incident-response-challenge/pkg/telemetry/otel"
	"github.com/lootarola/ai-incident-response-challenge/pkg/types"
	metricnoop "go.opentelemetry.io/otel/metric/noop"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	tracenoop "go.opentelemetry.io/otel/trace/noop"
)

func newTestService(t *testing.T) (*service.Service, *mock.OrderMock) {
	t.Helper()
	repo := mock.NewOrderMock()
	scope := &telemetryotel.Scope{
		Logger: telemetryotel.NewLogger(sdklog.NewLoggerProvider(), "test"),
		Tracer: tracenoop.NewTracerProvider().Tracer("test"),
		Meter:  metricnoop.NewMeterProvider().Meter("test"),
	}
	svc, err := service.NewService(repo, scope)
	if err != nil {
		t.Fatalf("NewService: %v", err)
	}
	return svc, repo
}

func seedOrder(t *testing.T, repo *mock.OrderMock) *types.Order {
	t.Helper()
	o := &types.Order{
		ID:         uuid.New().String(),
		CustomerID: "cust-001",
		Items: []types.Item{
			{ProductID: "prod-001", Category: "electronics", Quantity: 2, UnitPrice: 49.99},
		},
		Total:     99.98,
		Status:    types.StatusPending,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := repo.Create(context.Background(), o); err != nil {
		t.Fatalf("seed order: %v", err)
	}
	return o
}

func TestCreate(t *testing.T) {
	svc, _ := newTestService(t)

	req := orderdto.CreateRequest{
		CustomerID: "cust-002",
		Items:      []types.Item{{ProductID: "prod-001", Category: "books", Quantity: 1, UnitPrice: 12.99}},
	}
	o, err := svc.Create(context.Background(), req)
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if o.ID == "" {
		t.Error("expected non-empty ID")
	}
	if o.CustomerID != req.CustomerID {
		t.Errorf("customer_id: got %q, want %q", o.CustomerID, req.CustomerID)
	}
	if o.Status != types.StatusPending {
		t.Errorf("status: got %q, want %q", o.Status, types.StatusPending)
	}
	if o.Total != 12.99 {
		t.Errorf("total: got %f, want 12.99", o.Total)
	}
}

func TestGet(t *testing.T) {
	svc, repo := newTestService(t)
	seeded := seedOrder(t, repo)

	got, err := svc.Get(context.Background(), seeded.ID)
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if got.ID != seeded.ID {
		t.Errorf("id: got %q, want %q", got.ID, seeded.ID)
	}
}

func TestGet_NotFound(t *testing.T) {
	svc, _ := newTestService(t)

	_, err := svc.Get(context.Background(), "nonexistent")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !isOrderNotFound(err) {
		t.Errorf("expected ErrOrderNotFound, got %v", err)
	}
}

func TestUpdate(t *testing.T) {
	svc, repo := newTestService(t)
	seeded := seedOrder(t, repo)

	got, err := svc.Update(context.Background(), seeded.ID, orderdto.UpdateRequest{Status: types.StatusConfirmed})
	if err != nil {
		t.Fatalf("Update: %v", err)
	}
	if got.Status != types.StatusConfirmed {
		t.Errorf("status: got %q, want %q", got.Status, types.StatusConfirmed)
	}
}

func TestUpdate_InvalidTransition(t *testing.T) {
	svc, repo := newTestService(t)
	seeded := seedOrder(t, repo)

	_, err := svc.Update(context.Background(), seeded.ID, orderdto.UpdateRequest{Status: types.StatusDelivered})
	if err == nil {
		t.Fatal("expected error for invalid transition, got nil")
	}
}

func TestUpdate_NotFound(t *testing.T) {
	svc, _ := newTestService(t)

	_, err := svc.Update(context.Background(), "nonexistent", orderdto.UpdateRequest{Status: types.StatusConfirmed})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !isOrderNotFound(err) {
		t.Errorf("expected ErrOrderNotFound, got %v", err)
	}
}

func TestDelete(t *testing.T) {
	svc, repo := newTestService(t)
	seeded := seedOrder(t, repo)

	if err := svc.Delete(context.Background(), seeded.ID); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	if _, err := svc.Get(context.Background(), seeded.ID); err == nil {
		t.Error("expected order to be gone after delete")
	}
}

func TestDelete_NotFound(t *testing.T) {
	svc, _ := newTestService(t)

	err := svc.Delete(context.Background(), "nonexistent")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !isOrderNotFound(err) {
		t.Errorf("expected ErrOrderNotFound, got %v", err)
	}
}

func TestReport(t *testing.T) {
	svc, repo := newTestService(t)
	seedOrder(t, repo)

	resp, err := svc.Report(context.Background())
	if err != nil {
		t.Fatalf("Report: %v", err)
	}
	if resp == nil {
		t.Fatal("expected non-nil response")
	}
	if resp.TotalCount != len(resp.Entries) {
		t.Errorf("TotalCount %d != len(Entries) %d", resp.TotalCount, len(resp.Entries))
	}
}

func TestNotify_OrderNotFound(t *testing.T) {
	svc, _ := newTestService(t)

	err := svc.Notify(context.Background(), "nonexistent", "confirmed")
	if err == nil {
		t.Fatal("expected error for missing order, got nil")
	}
	if !isOrderNotFound(err) {
		t.Errorf("expected ErrOrderNotFound, got %v", err)
	}
}

func TestNotify_Succeeds(t *testing.T) {
	svc, repo := newTestService(t)
	seeded := seedOrder(t, repo)

	if err := svc.Notify(context.Background(), seeded.ID, "confirmed"); err != nil {
		t.Fatalf("Notify: %v", err)
	}
}

func isOrderNotFound(err error) bool {
	return containsError(err, types.ErrOrderNotFound)
}

func containsError(err, target error) bool {
	for err != nil {
		if err == target {
			return true
		}
		type unwrapper interface{ Unwrap() error }
		if u, ok := err.(unwrapper); ok {
			err = u.Unwrap()
		} else {
			return false
		}
	}
	return false
}
