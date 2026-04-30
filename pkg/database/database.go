package database

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/lootarola/ai-incident-response-challenge/pkg/types"
)

var (
	ErrNotFound = errors.New("not found")
	ErrConflict = errors.New("conflict")
)

type OrderRepository interface {
	Create(ctx context.Context, order *types.Order) error
	Get(ctx context.Context, id string) (*types.Order, error)
	Update(ctx context.Context, order *types.Order) error
	Delete(ctx context.Context, id string) error
	ListSince(ctx context.Context, since time.Time) ([]*types.Order, error)
}

type CatalogRepository interface {
	Create(ctx context.Context, p *types.Product) error
	AggregatedGetByID(ctx context.Context, id string) (*types.Product, error)
	Update(ctx context.Context, p *types.Product) error
	Delete(ctx context.Context, id string) error
	Search(ctx context.Context, filter map[string]interface{}) ([]*types.Product, error)
}

type OrderFactory func(ctx context.Context) (OrderRepository, func(context.Context) error, error)
type CatalogFactory func(ctx context.Context) (CatalogRepository, func(context.Context) error, error)

var (
	orderMu        sync.RWMutex
	orderFactories = map[string]OrderFactory{}

	catalogMu        sync.RWMutex
	catalogFactories = map[string]CatalogFactory{}
)

func RegisterOrder(name string, f OrderFactory) {
	orderMu.Lock()
	defer orderMu.Unlock()
	orderFactories[name] = f
}

func RegisterCatalog(name string, f CatalogFactory) {
	catalogMu.Lock()
	defer catalogMu.Unlock()
	catalogFactories[name] = f
}

func NewOrder(ctx context.Context) (OrderRepository, func(context.Context) error, error) {
	backend := os.Getenv("ORDER_DB_BACKEND")
	orderMu.RLock()
	f, ok := orderFactories[backend]
	orderMu.RUnlock()
	if !ok {
		return nil, nil, fmt.Errorf("ORDER_DB_BACKEND %q not registered", backend)
	}
	return f(ctx)
}

func NewCatalog(ctx context.Context) (CatalogRepository, func(context.Context) error, error) {
	backend := os.Getenv("CATALOG_DB_BACKEND")
	catalogMu.RLock()
	f, ok := catalogFactories[backend]
	catalogMu.RUnlock()
	if !ok {
		return nil, nil, fmt.Errorf("CATALOG_DB_BACKEND %q not registered", backend)
	}
	return f(ctx)
}
