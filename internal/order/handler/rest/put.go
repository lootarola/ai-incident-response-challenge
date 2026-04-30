package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	orderdto "github.com/lootarola/ai-incident-response-challenge/internal/order/handler/rest/dto"
)

func (h *Handler) Update(c *gin.Context) {
	id := c.Param("id")

	var req orderdto.UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	o, err := h.svc.Update(c.Request.Context(), id, req)
	if err != nil {
		c.JSON(httpStatus(err), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, o)
}
