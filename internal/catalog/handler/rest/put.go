package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	catalogdto "github.com/lootarola/ai-incident-response-challenge/internal/catalog/handler/rest/dto"
)

func (h *Handler) Update(c *gin.Context) {
	id := c.Param("id")

	var req catalogdto.UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	p, err := h.svc.Update(c.Request.Context(), id, req)
	if err != nil {
		c.JSON(httpStatus(err), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, p)
}
