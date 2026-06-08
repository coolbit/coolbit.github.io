package handlers

import (
	"blog/store"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type PostHandler struct {
	s *store.Store
}

func NewPostHandler(s *store.Store) *PostHandler {
	return &PostHandler{s: s}
}

type postInput struct {
	Title    string `json:"title" binding:"required"`
	Summary  string `json:"summary"`
	Content  string `json:"content"`
	Category string `json:"category"`
}

func (h *PostHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	posts, err := h.s.All()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if cat := c.Query("category"); cat != "" {
		filtered := posts[:0]
		for _, p := range posts {
			if p.Category == cat {
				filtered = append(filtered, p)
			}
		}
		posts = filtered
	}

	if kw := strings.ToLower(c.Query("q")); kw != "" {
		filtered := posts[:0]
		for _, p := range posts {
			if strings.Contains(strings.ToLower(p.Title), kw) ||
				strings.Contains(strings.ToLower(p.Summary), kw) ||
				strings.Contains(strings.ToLower(p.Content), kw) {
				filtered = append(filtered, p)
			}
		}
		posts = filtered
	}

	total := len(posts)
	start := (page - 1) * pageSize
	if start > total {
		start = total
	}
	end := start + pageSize
	if end > total {
		end = total
	}

	c.JSON(http.StatusOK, gin.H{
		"data":      posts[start:end],
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func (h *PostHandler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	p, err := h.s.Get(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, p)
}

func (h *PostHandler) Create(c *gin.Context) {
	var input postInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	p, err := h.s.Create(input.Title, input.Summary, input.Content, input.Category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, p)
}

func (h *PostHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var input postInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	p, err := h.s.Update(uint(id), input.Title, input.Summary, input.Content, input.Category)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, p)
}

func (h *PostHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.s.Delete(uint(id)); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.Status(http.StatusNoContent)
}
