package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) Delete(c *gin.Context) {
	id := c.Param("id")

	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		c.JSON(httpStatus(err), gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
