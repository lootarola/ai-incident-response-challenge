package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	catalogdto "github.com/lootarola/ai-incident-response-challenge/internal/catalog/handler/rest/dto"
	"github.com/lootarola/ai-incident-response-challenge/internal/catalog/service"
	"github.com/lootarola/ai-incident-response-challenge/pkg/database/mock"
	telemetryotel "github.com/lootarola/ai-incident-response-challenge/pkg/telemetry/otel"
	"github.com/lootarola/ai-incident-response-challenge/pkg/types"
	metricnoop "go.opentelemetry.io/otel/metric/noop"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	tracenoop "go.opentelemetry.io/otel/trace/noop"
)

func newTestService(t *testing.T) (*service.Service, *mock.CatalogMock) {
	t.Helper()
	repo := mock.NewCatalogMock()
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

func seedProduct(t *testing.T, repo *mock.CatalogMock) *types.Product {
	t.Helper()
	p := &types.Product{
		ID:          uuid.New().String(),
		Name:        "Widget",
		Description: "A standard widget",
		Category:    "electronics",
		Price:       49.99,
		Inventory:   100,
		IsInternal:  false,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	if err := repo.Create(context.Background(), p); err != nil {
		t.Fatalf("seed product: %v", err)
	}
	return p
}

func TestCreate(t *testing.T) {
	svc, _ := newTestService(t)

	req := catalogdto.CreateRequest{
		Name:     "Gadget",
		Category: "electronics",
		Price:    99.99,
	}
	p, err := svc.Create(context.Background(), req)
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if p.ID == "" {
		t.Error("expected non-empty ID")
	}
	if p.Name != req.Name {
		t.Errorf("name: got %q, want %q", p.Name, req.Name)
	}
	if p.Price != req.Price {
		t.Errorf("price: got %f, want %f", p.Price, req.Price)
	}
}

func TestGetByID(t *testing.T) {
	svc, repo := newTestService(t)
	seeded := seedProduct(t, repo)

	got, err := svc.GetByID(context.Background(), seeded.ID)
	if err != nil {
		t.Fatalf("GetByID: %v", err)
	}
	if got.ID != seeded.ID {
		t.Errorf("id: got %q, want %q", got.ID, seeded.ID)
	}
}

func TestGetByID_NotFound(t *testing.T) {
	svc, _ := newTestService(t)

	_, err := svc.GetByID(context.Background(), "nonexistent")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !isProductNotFound(err) {
		t.Errorf("expected ErrProductNotFound, got %v", err)
	}
}

func TestUpdate(t *testing.T) {
	svc, repo := newTestService(t)
	seeded := seedProduct(t, repo)

	got, err := svc.Update(context.Background(), seeded.ID, catalogdto.UpdateRequest{Name: "Updated Widget", Price: 59.99})
	if err != nil {
		t.Fatalf("Update: %v", err)
	}
	if got.Name != "Updated Widget" {
		t.Errorf("name: got %q, want %q", got.Name, "Updated Widget")
	}
	if got.Price != 59.99 {
		t.Errorf("price: got %f, want 59.99", got.Price)
	}
}

func TestUpdate_NotFound(t *testing.T) {
	svc, _ := newTestService(t)

	_, err := svc.Update(context.Background(), "nonexistent", catalogdto.UpdateRequest{Name: "X"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !isProductNotFound(err) {
		t.Errorf("expected ErrProductNotFound, got %v", err)
	}
}

func TestDelete(t *testing.T) {
	svc, repo := newTestService(t)
	seeded := seedProduct(t, repo)

	if err := svc.Delete(context.Background(), seeded.ID); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	if _, err := svc.GetByID(context.Background(), seeded.ID); err == nil {
		t.Error("expected product to be gone after delete")
	}
}

func TestDelete_NotFound(t *testing.T) {
	svc, _ := newTestService(t)

	err := svc.Delete(context.Background(), "nonexistent")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !isProductNotFound(err) {
		t.Errorf("expected ErrProductNotFound, got %v", err)
	}
}

func TestSearch_SafePath(t *testing.T) {
	svc, repo := newTestService(t)
	seedProduct(t, repo)

	results, err := svc.Search(context.Background(), "widget")
	if err != nil {
		t.Fatalf("Search: %v", err)
	}
	if results == nil {
		t.Error("expected non-nil results slice")
	}
}

func TestSearch_EmptyQuery(t *testing.T) {
	svc, repo := newTestService(t)
	seedProduct(t, repo)

	results, err := svc.Search(context.Background(), "")
	if err != nil {
		t.Fatalf("Search with empty query: %v", err)
	}
	if results == nil {
		t.Error("expected non-nil results slice")
	}
}

func isProductNotFound(err error) bool {
	return containsError(err, types.ErrProductNotFound)
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
