package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) Report(c *gin.Context) {
	report, err := h.svc.Report(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, report)
}
