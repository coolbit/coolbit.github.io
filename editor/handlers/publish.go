package handlers

import (
	"blog/publisher"
	"blog/store"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

const postsPerPage = 10

type PublishHandler struct {
	s *store.Store
}

func NewPublishHandler(s *store.Store) *PublishHandler {
	return &PublishHandler{s: s}
}

func (h *PublishHandler) Publish(c *gin.Context) {
	outDir := os.Getenv("PUBLISH_DIR")
	if outDir == "" {
		outDir = ".."
	}
	siteURL := os.Getenv("SITE_URL")
	if siteURL == "" {
		siteURL = "https://coolbit.github.io"
	}

	if err := h.s.RerenderAll(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "rerender: " + err.Error()})
		return
	}

	posts, err := h.s.All()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	type job struct{ path, content string }
	jobs := []job{
		{"style.css", publisher.StyleCSS},
		{"sitemap.xml", publisher.GenerateSitemap(posts, siteURL)},
	}

	for _, pg := range publisher.PaginateIndex(posts, postsPerPage) {
		html, err := publisher.RenderIndex(pg.Data)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "render: " + err.Error()})
			return
		}
		jobs = append(jobs, job{pg.Path, html})
	}

	var errs []string
	for _, j := range jobs {
		full := filepath.Join(outDir, j.path)
		if err := os.MkdirAll(filepath.Dir(full), 0755); err != nil {
			errs = append(errs, full+": "+err.Error())
			continue
		}
		if err := os.WriteFile(full, []byte(j.content), 0644); err != nil {
			errs = append(errs, full+": "+err.Error())
		}
	}

	if len(errs) > 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": errs})
		return
	}

	abs, _ := filepath.Abs(outDir)

	cmd := exec.Command("npx", "pagefind", "--site", abs, "--output-subdir", "_pagefind", "--glob", "posts/*.html")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Printf("pagefind: %v", err)
	}

	c.JSON(http.StatusOK, gin.H{"dir": abs, "count": len(posts)})
}
