package types

import (
	"errors"
	"time"
)

type Product struct {
	ID          string            `bson:"_id" json:"id"`
	Name        string            `bson:"name" json:"name"`
	Description string            `bson:"description" json:"description"`
	Category    string            `bson:"category" json:"category"`
	Price       float64           `bson:"price" json:"price"`
	Inventory   int               `bson:"inventory" json:"inventory"`
	IsInternal  bool              `bson:"is_internal" json:"is_internal"`
	Reviews     []Review          `bson:"reviews" json:"reviews"`
	Specs       map[string]string `bson:"specs" json:"specs"`
	CreatedAt   time.Time         `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time         `bson:"updated_at" json:"updated_at"`
}

type Review struct {
	Author    string    `bson:"author" json:"author"`
	Body      string    `bson:"body" json:"body"`
	Rating    int       `bson:"rating" json:"rating"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
}

var (
	ErrProductNotFound = errors.New("product not found")
	ErrInvalidProduct  = errors.New("invalid product")
	ErrInvalidSearch   = errors.New("invalid search query")
)
