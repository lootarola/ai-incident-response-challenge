package mock

import (
	"context"
	"sync"
	"time"

	"github.com/lootarola/ai-incident-response-challenge/pkg/database"
	"github.com/lootarola/ai-incident-response-challenge/pkg/types"
)

type OrderMock struct {
	mu     sync.RWMutex
	orders map[string]*types.Order
}

func NewOrderMock() *OrderMock {
	return &OrderMock{orders: make(map[string]*types.Order)}
}

func (m *OrderMock) Create(_ context.Context, order *types.Order) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.orders[order.ID]; ok {
		return database.ErrConflict
	}
	cp := *order
	m.orders[order.ID] = &cp
	return nil
}

func (m *OrderMock) Get(_ context.Context, id string) (*types.Order, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	o, ok := m.orders[id]
	if !ok {
		return nil, database.ErrNotFound
	}
	cp := *o
	return &cp, nil
}

func (m *OrderMock) Update(_ context.Context, order *types.Order) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.orders[order.ID]; !ok {
		return database.ErrNotFound
	}
	cp := *order
	m.orders[order.ID] = &cp
	return nil
}

func (m *OrderMock) Delete(_ context.Context, id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.orders[id]; !ok {
		return database.ErrNotFound
	}
	delete(m.orders, id)
	return nil
}

func (m *OrderMock) ListSince(_ context.Context, since time.Time) ([]*types.Order, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var out []*types.Order
	for _, o := range m.orders {
		if !o.CreatedAt.Before(since) {
			cp := *o
			out = append(out, &cp)
		}
	}
	return out, nil
}
