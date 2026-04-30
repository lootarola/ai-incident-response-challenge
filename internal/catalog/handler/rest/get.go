package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	catalogdto "github.com/lootarola/ai-incident-response-challenge/internal/catalog/handler/rest/dto"
	"github.com/lootarola/ai-incident-response-challenge/pkg/types"
)

func (h *Handler) GetByID(c *gin.Context) {
	id := c.Param("id")

	p, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(httpStatus(err), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, p)
}

func (h *Handler) Search(c *gin.Context) {
	raw := c.Query("search")

	products, err := h.svc.Search(c.Request.Context(), raw)
	if err != nil {
		c.JSON(httpStatus(err), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, catalogdto.SearchResponse{
		Products: func() []types.Product {
			out := make([]types.Product, 0, len(products))
			for _, p := range products {
				out = append(out, *p)
			}
			return out
		}(),
		Count: len(products),
	})
}
