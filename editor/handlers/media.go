package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

type MediaListHandler struct {
	mediaDir string
}

func NewMediaListHandler(mediaDir string) *MediaListHandler {
	return &MediaListHandler{mediaDir: mediaDir}
}

func (h *MediaListHandler) List(c *gin.Context) {
	postID := c.Param("id")
	dir := filepath.Join(h.mediaDir, postID)
	entries, err := os.ReadDir(dir)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"files": []string{}})
		return
	}
	files := make([]string, 0, len(entries))
	for _, e := range entries {
		if !e.IsDir() {
			files = append(files, fmt.Sprintf("/media/%s/%s", postID, e.Name()))
		}
	}
	c.JSON(http.StatusOK, gin.H{"files": files})
}
