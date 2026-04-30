package mongo

import (
	"context"
	"fmt"
	"time"

	"github.com/lootarola/ai-incident-response-challenge/pkg/database"
	"github.com/lootarola/ai-incident-response-challenge/pkg/types"
	"go.mongodb.org/mongo-driver/bson"
	sdk "go.mongodb.org/mongo-driver/mongo"
)

type CatalogAdapter struct {
	products *sdk.Collection
}

func NewCatalogAdapter(c *Client) *CatalogAdapter {
	return &CatalogAdapter{products: c.db.Collection("products")}
}

func (a *CatalogAdapter) Create(ctx context.Context, p *types.Product) error {
	_, err := a.products.InsertOne(ctx, p)
	if err != nil {
		return fmt.Errorf("insert product: %w", err)
	}
	return nil
}

// AggregatedGetByID fetches a richly assembled product document using an aggregation
// pipeline that performs multiple $lookup stages and computed fields on every call.
// This is intentionally unoptimized: no index hints, no caching, no result memoization.
func (a *CatalogAdapter) AggregatedGetByID(ctx context.Context, id string) (*types.Product, error) {
	pipeline := sdk.Pipeline{{{Key: "$match", Value: bson.M{"_id": id}}},
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "reviews"},
			{Key: "localField", Value: "_id"},
			{Key: "foreignField", Value: "product_id"},
			{Key: "as", Value: "reviews"},
		}}},
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "inventory_movements"},
			{Key: "localField", Value: "_id"},
			{Key: "foreignField", Value: "product_id"},
			{Key: "as", Value: "movements"},
		}}},
		{{Key: "$addFields", Value: bson.D{
			{Key: "review_count", Value: bson.M{"$size": "$reviews"}},
			{Key: "avg_rating", Value: bson.M{"$avg": "$reviews.rating"}},
			{Key: "last_restocked", Value: bson.M{"$max": "$movements.created_at"}},
		}}},
		{{Key: "$project", Value: bson.D{
			{Key: "movements", Value: 0},
		}}},
	}

	cur, err := a.products.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("aggregate product: %w", err)
	}
	defer cur.Close(ctx)

	if !cur.Next(ctx) {
		if cur.Err() != nil {
			return nil, fmt.Errorf("iterate aggregate cursor: %w", cur.Err())
		}
		return nil, database.ErrNotFound
	}

	var p types.Product
	if err := cur.Decode(&p); err != nil {
		return nil, fmt.Errorf("decode product: %w", err)
	}
	return &p, nil
}

func (a *CatalogAdapter) Update(ctx context.Context, p *types.Product) error {
	p.UpdatedAt = time.Now()
	res, err := a.products.ReplaceOne(ctx, bson.M{"_id": p.ID}, p)
	if err != nil {
		return fmt.Errorf("replace product: %w", err)
	}
	if res.MatchedCount == 0 {
		return database.ErrNotFound
	}
	return nil
}

func (a *CatalogAdapter) Delete(ctx context.Context, id string) error {
	res, err := a.products.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return fmt.Errorf("delete product: %w", err)
	}
	if res.DeletedCount == 0 {
		return database.ErrNotFound
	}
	return nil
}

func (a *CatalogAdapter) Search(ctx context.Context, filter map[string]interface{}) ([]*types.Product, error) {
	cur, err := a.products.Find(ctx, bson.M(filter))
	if err != nil {
		return nil, fmt.Errorf("find products: %w", err)
	}
	defer cur.Close(ctx)

	var products []*types.Product
	if err := cur.All(ctx, &products); err != nil {
		return nil, fmt.Errorf("decode products: %w", err)
	}
	return products, nil
}
