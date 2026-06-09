package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

type UploadHandler struct {
	mediaDir string
}

func NewUploadHandler(mediaDir string) *UploadHandler {
	return &UploadHandler{mediaDir: mediaDir}
}

func (h *UploadHandler) Upload(c *gin.Context) {
	postID := c.Param("id")

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	filename := strings.ReplaceAll(filepath.Base(file.Filename), " ", "-")

	dir := filepath.Join(h.mediaDir, postID)
	if err := os.MkdirAll(dir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := c.SaveUploadedFile(file, filepath.Join(dir, filename)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"url": fmt.Sprintf("/media/%s/%s", postID, filename)})
}
