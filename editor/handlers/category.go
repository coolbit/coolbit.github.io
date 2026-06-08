package handlers

import (
	"blog/store"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	s *store.Store
}

func NewCategoryHandler(s *store.Store) *CategoryHandler {
	return &CategoryHandler{s: s}
}

func (h *CategoryHandler) List(c *gin.Context) {
	c.JSON(http.StatusOK, h.s.Categories())
}
