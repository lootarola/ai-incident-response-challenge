package dto

import "github.com/lootarola/ai-incident-response-challenge/pkg/types"

type CreateRequest struct {
	CustomerID string       `json:"customer_id" binding:"required"`
	Items      []types.Item `json:"items" binding:"required,min=1"`
}

type UpdateRequest struct {
	Status types.OrderStatus `json:"status" binding:"required"`
}

type NotifyRequest struct {
	Event string `json:"event" binding:"required"`
}

type ReportResponse struct {
	Entries    []types.Report `json:"entries"`
	TotalCount int            `json:"total_count"`
}
