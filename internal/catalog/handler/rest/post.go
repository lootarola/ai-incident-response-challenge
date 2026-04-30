package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	catalogdto "github.com/lootarola/ai-incident-response-challenge/internal/catalog/handler/rest/dto"
)

func (h *Handler) Create(c *gin.Context) {
	var req catalogdto.CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	p, err := h.svc.Create(c.Request.Context(), req)
	if err != nil {
		c.JSON(httpStatus(err), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, p)
}
