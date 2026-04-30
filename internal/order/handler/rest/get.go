package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetByID(c *gin.Context) {
	id := c.Param("id")

	o, err := h.svc.Get(c.Request.Context(), id)
	if err != nil {
		c.JSON(httpStatus(err), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, o)
}
