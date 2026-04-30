package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	orderdto "github.com/lootarola/ai-incident-response-challenge/internal/order/handler/rest/dto"
)

func (h *Handler) Notify(c *gin.Context) {
	id := c.Param("id")

	var req orderdto.NotifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.svc.Notify(c.Request.Context(), id, req.Event); err != nil {
		c.JSON(httpStatus(err), gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusAccepted)
}
