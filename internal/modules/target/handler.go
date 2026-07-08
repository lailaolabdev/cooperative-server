package target

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct{ service *Service }

func NewHandler(service *Service) *Handler { return &Handler{service: service} }

func (h *Handler) List(c *gin.Context) {
	items, summary, err := h.service.List(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to load cooperative targets"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": items, "summary": summary})
}
