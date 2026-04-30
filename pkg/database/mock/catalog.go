package mock

import (
	"context"
	"sync"

	"github.com/lootarola/ai-incident-response-challenge/pkg/database"
	"github.com/lootarola/ai-incident-response-challenge/pkg/types"
)

type CatalogMock struct {
	mu       sync.RWMutex
	products map[string]*types.Product
}

func NewCatalogMock() *CatalogMock {
	return &CatalogMock{products: make(map[string]*types.Product)}
}

func (m *CatalogMock) Create(_ context.Context, p *types.Product) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.products[p.ID]; ok {
		return database.ErrConflict
	}
	cp := *p
	m.products[p.ID] = &cp
	return nil
}

func (m *CatalogMock) AggregatedGetByID(_ context.Context, id string) (*types.Product, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	p, ok := m.products[id]
	if !ok {
		return nil, database.ErrNotFound
	}
	cp := *p
	return &cp, nil
}

func (m *CatalogMock) Update(_ context.Context, p *types.Product) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.products[p.ID]; !ok {
		return database.ErrNotFound
	}
	cp := *p
	m.products[p.ID] = &cp
	return nil
}

func (m *CatalogMock) Delete(_ context.Context, id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.products[id]; !ok {
		return database.ErrNotFound
	}
	delete(m.products, id)
	return nil
}

func (m *CatalogMock) Search(_ context.Context, _ map[string]interface{}) ([]*types.Product, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]*types.Product, 0, len(m.products))
	for _, p := range m.products {
		cp := *p
		out = append(out, &cp)
	}
	return out, nil
}
