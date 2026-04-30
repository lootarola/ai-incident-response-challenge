package rest

import (
	"errors"

	"github.com/gin-gonic/gin"
	catalogservice "github.com/lootarola/ai-incident-response-challenge/internal/catalog/service"
	"github.com/lootarola/ai-incident-response-challenge/pkg/types"
)

type Handler struct {
	svc *catalogservice.Service
}

func NewHandler(svc *catalogservice.Service) (*Handler, error) {
	if svc == nil {
		return nil, errors.New("service is required")
	}
	return &Handler{svc: svc}, nil
}

func (h *Handler) RegisterRoutes(r gin.IRouter) {
	r.GET("/catalog/products", h.Search)
	r.GET("/catalog/products/:id", h.GetByID)
	r.POST("/catalog/products", h.Create)
	r.PUT("/catalog/products/:id", h.Update)
	r.DELETE("/catalog/products/:id", h.Delete)
}

func httpStatus(err error) int {
	switch {
	case errors.Is(err, types.ErrProductNotFound):
		return 404
	case errors.Is(err, types.ErrInvalidProduct),
		errors.Is(err, types.ErrInvalidSearch):
		return 400
	default:
		return 500
	}
}
