package dto

import "github.com/lootarola/ai-incident-response-challenge/pkg/types"

type CreateRequest struct {
	Name        string            `json:"name" binding:"required"`
	Description string            `json:"description"`
	Category    string            `json:"category" binding:"required"`
	Price       float64           `json:"price" binding:"required,gt=0"`
	Inventory   int               `json:"inventory"`
	IsInternal  bool              `json:"is_internal"`
	Specs       map[string]string `json:"specs"`
}

type UpdateRequest struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Category    string            `json:"category"`
	Price       float64           `json:"price"`
	Inventory   int               `json:"inventory"`
	Specs       map[string]string `json:"specs"`
}

type SearchResponse struct {
	Products []types.Product `json:"products"`
	Count    int             `json:"count"`
}
