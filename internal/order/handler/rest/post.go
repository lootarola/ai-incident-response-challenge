package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	orderdto "github.com/lootarola/ai-incident-response-challenge/internal/order/handler/rest/dto"
)

func (h *Handler) Create(c *gin.Context) {
	var req orderdto.CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	o, err := h.svc.Create(c.Request.Context(), req)
	if err != nil {
		c.JSON(httpStatus(err), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, o)
}
