package mongo

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/lootarola/ai-incident-response-challenge/pkg/database"
	"github.com/lootarola/ai-incident-response-challenge/pkg/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	sdk "go.mongodb.org/mongo-driver/mongo"
)

type OrderAdapter struct {
	coll *sdk.Collection
}

func NewOrderAdapter(c *Client) *OrderAdapter {
	return &OrderAdapter{coll: c.db.Collection("orders")}
}

func (a *OrderAdapter) Create(ctx context.Context, order *types.Order) error {
	_, err := a.coll.InsertOne(ctx, order)
	if err != nil {
		return fmt.Errorf("insert order: %w", err)
	}
	return nil
}

func (a *OrderAdapter) Get(ctx context.Context, id string) (*types.Order, error) {
	var order types.Order
	err := a.coll.FindOne(ctx, bson.M{"_id": id}).Decode(&order)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, database.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("find order: %w", err)
	}
	return &order, nil
}

func (a *OrderAdapter) Update(ctx context.Context, order *types.Order) error {
	order.UpdatedAt = time.Now()
	res, err := a.coll.ReplaceOne(ctx, bson.M{"_id": order.ID}, order)
	if err != nil {
		return fmt.Errorf("replace order: %w", err)
	}
	if res.MatchedCount == 0 {
		return database.ErrNotFound
	}
	return nil
}

func (a *OrderAdapter) Delete(ctx context.Context, id string) error {
	res, err := a.coll.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return fmt.Errorf("delete order: %w", err)
	}
	if res.DeletedCount == 0 {
		return database.ErrNotFound
	}
	return nil
}

func (a *OrderAdapter) ListSince(ctx context.Context, since time.Time) ([]*types.Order, error) {
	cur, err := a.coll.Find(ctx, bson.M{"created_at": bson.M{"$gte": since}})
	if err != nil {
		return nil, fmt.Errorf("find orders: %w", err)
	}
	defer cur.Close(ctx)

	var orders []*types.Order
	if err := cur.All(ctx, &orders); err != nil {
		return nil, fmt.Errorf("decode orders: %w", err)
	}
	return orders, nil
}
