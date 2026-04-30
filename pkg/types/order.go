package types

import (
	"errors"
	"time"
)

type Order struct {
	ID         string    `bson:"_id" json:"id"`
	CustomerID string    `bson:"customer_id" json:"customer_id"`
	Items      []Item    `bson:"items" json:"items"`
	Total      float64   `bson:"total" json:"total"`
	Status     OrderStatus `bson:"status" json:"status"`
	CreatedAt  time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time `bson:"updated_at" json:"updated_at"`
}

type Item struct {
	ProductID string  `bson:"product_id" json:"product_id"`
	Category  string  `bson:"category" json:"category"`
	Quantity  int     `bson:"quantity" json:"quantity"`
	UnitPrice float64 `bson:"unit_price" json:"unit_price"`
}

type Report struct {
	CustomerID string    `json:"customer_id"`
	Category   string    `json:"category"`
	TotalSpend float64   `json:"total_spend"`
	OrderCount int       `json:"order_count"`
	Detail     string    `json:"detail"`
	ComputedAt time.Time `json:"computed_at"`
}

var (
	ErrOrderNotFound        = errors.New("order not found")
	ErrInvalidOrder         = errors.New("invalid order")
	ErrInventoryUnavailable = errors.New("inventory service unavailable")
	ErrInvalidStatus        = errors.New("invalid order status")
)

type OrderStatus string

const (
	StatusPending   OrderStatus = "pending"
	StatusConfirmed OrderStatus = "confirmed"
	StatusShipped   OrderStatus = "shipped"
	StatusDelivered OrderStatus = "delivered"
	StatusCancelled OrderStatus = "cancelled"
)

type NotificationChannel string

const (
	ChannelEmail   NotificationChannel = "email"
	ChannelPush    NotificationChannel = "push"
	ChannelSMS     NotificationChannel = "sms"
	ChannelWebhook NotificationChannel = "webhook"
)

var NotificationChannels = []NotificationChannel{ChannelEmail, ChannelPush, ChannelSMS, ChannelWebhook}
