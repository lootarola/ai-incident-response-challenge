package mongo

import (
	"context"
	"fmt"

	sdk "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.opentelemetry.io/contrib/instrumentation/go.mongodb.org/mongo-driver/mongo/otelmongo"
)

type Client struct {
	client *sdk.Client
	db     *sdk.Database
}

func NewClient(ctx context.Context, uri, dbName string) (*Client, error) {
	opts := options.Client().
		ApplyURI(uri).
		SetMonitor(otelmongo.NewMonitor())

	c, err := sdk.Connect(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("connect to sdkdb: %w", err)
	}

	if err := c.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("ping sdkdb: %w", err)
	}

	return &Client{client: c, db: c.Database(dbName)}, nil
}

func (c *Client) Database() *sdk.Database {
	return c.db
}

func (c *Client) Disconnect(ctx context.Context) error {
	return c.client.Disconnect(ctx)
}
