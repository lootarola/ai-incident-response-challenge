package mock

import (
	"context"

	"github.com/lootarola/ai-incident-response-challenge/pkg/database"
)

func init() {
	database.RegisterOrder("mock", newOrder)
	database.RegisterCatalog("mock", newCatalog)
}

func newOrder(_ context.Context) (database.OrderRepository, func(context.Context) error, error) {
	return NewOrderMock(), func(_ context.Context) error { return nil }, nil
}

func newCatalog(_ context.Context) (database.CatalogRepository, func(context.Context) error, error) {
	return NewCatalogMock(), func(_ context.Context) error { return nil }, nil
}
