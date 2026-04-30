package rest

import (
	"errors"

	"github.com/gin-gonic/gin"
	orderservice "github.com/lootarola/ai-incident-response-challenge/internal/order/service"
	"github.com/lootarola/ai-incident-response-challenge/pkg/types"
)

type Handler struct {
	svc *orderservice.Service
}

func NewHandler(svc *orderservice.Service) (*Handler, error) {
	if svc == nil {
		return nil, errors.New("service is required")
	}
	return &Handler{svc: svc}, nil
}

func (h *Handler) RegisterRoutes(r gin.IRouter) {
	r.GET("/orders/report", h.Report)
	r.GET("/orders/:id", h.GetByID)
	r.POST("/orders", h.Create)
	r.POST("/orders/:id/notify", h.Notify)
	r.PUT("/orders/:id", h.Update)
	r.DELETE("/orders/:id", h.Delete)
}

func httpStatus(err error) int {
	switch {
	case errors.Is(err, types.ErrOrderNotFound):
		return 404
	case errors.Is(err, types.ErrInventoryUnavailable):
		return 503
	case errors.Is(err, types.ErrInvalidOrder),
		errors.Is(err, types.ErrInvalidStatus):
		return 400
	default:
		return 500
	}
}
