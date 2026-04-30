package mongo

import (
	"context"
	"os"
	"sync"

	"github.com/lootarola/ai-incident-response-challenge/pkg/database"
)

func init() {
	database.RegisterOrder("mongo", newOrder)
	database.RegisterCatalog("mongo", newCatalog)
}

var (
	sharedMu     sync.Mutex
	sharedClient *Client
	sharedRefs   int
)

func acquireClient(ctx context.Context) (*Client, error) {
	sharedMu.Lock()
	defer sharedMu.Unlock()
	if sharedClient == nil {
		c, err := NewClient(ctx, os.Getenv("MONGO_URI"), os.Getenv("MONGO_DB"))
		if err != nil {
			return nil, err
		}
		sharedClient = c
	}
	sharedRefs++
	return sharedClient, nil
}

func releaseClient(ctx context.Context) error {
	sharedMu.Lock()
	defer sharedMu.Unlock()
	sharedRefs--
	if sharedRefs == 0 {
		err := sharedClient.Disconnect(ctx)
		sharedClient = nil
		return err
	}
	return nil
}

func newOrder(ctx context.Context) (database.OrderRepository, func(context.Context) error, error) {
	c, err := acquireClient(ctx)
	if err != nil {
		return nil, nil, err
	}
	return NewOrderAdapter(c), releaseClient, nil
}

func newCatalog(ctx context.Context) (database.CatalogRepository, func(context.Context) error, error) {
	c, err := acquireClient(ctx)
	if err != nil {
		return nil, nil, err
	}
	return NewCatalogAdapter(c), releaseClient, nil
}
