package cooperative

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

type Handler struct{ service *Service }

func NewHandler(s *Service) *Handler { return &Handler{service: s} }
func (h *Handler) List(c *gin.Context) {
	items, e := h.service.List(c, c.Query("type"), c.Query("province"))
	if e != nil {
		fail(c, e)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": items})
}
func (h *Handler) Get(c *gin.Context) {
	item, e := h.service.Get(c, c.Param("id"))
	if e != nil {
		fail(c, e)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": item})
}
func (h *Handler) Create(c *gin.Context) {
	var q UpsertRequest
	if c.ShouldBindJSON(&q) != nil {
		c.JSON(400, gin.H{"error": "invalid cooperative data"})
		return
	}
	item, e := h.service.Create(c, q)
	if e != nil {
		fail(c, e)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": item})
}
func (h *Handler) Update(c *gin.Context) {
	var q UpsertRequest
	if c.ShouldBindJSON(&q) != nil {
		c.JSON(400, gin.H{"error": "invalid cooperative data"})
		return
	}
	item, e := h.service.Update(c, c.Param("id"), q)
	if e != nil {
		fail(c, e)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": item})
}
func (h *Handler) Delete(c *gin.Context) {
	if e := h.service.Delete(c, c.Param("id")); e != nil {
		fail(c, e)
		return
	}
	c.Status(http.StatusNoContent)
}
func fail(c *gin.Context, e error) {
	if errors.Is(e, mongo.ErrNoDocuments) {
		c.JSON(404, gin.H{"error": "cooperative not found"})
		return
	}
	if errors.Is(e, ErrInvalidType) {
		c.JSON(422, gin.H{"error": e.Error()})
		return
	}
	c.JSON(500, gin.H{"error": "internal server error"})
}
