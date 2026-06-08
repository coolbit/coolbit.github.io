package main

import (
	"blog/handlers"
	"blog/store"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	buildFrontend()

	postsDir := os.Getenv("POSTS_DIR")
	if postsDir == "" {
		postsDir = "../posts"
	}
	s, err := store.New(postsDir)
	if err != nil {
		log.Fatalf("store: %v", err)
	}

	r := gin.Default()
	r.Use(cors.Default())

	ph := handlers.NewPostHandler(s)
	ch := handlers.NewCategoryHandler(s)
	pub := handlers.NewPublishHandler(s)

	posts := r.Group("/api/posts")
	{
		posts.GET("", ph.List)
		posts.GET("/:id", ph.Get)
		posts.POST("", ph.Create)
		posts.PUT("/:id", ph.Update)
		posts.DELETE("/:id", ph.Delete)
	}

	r.GET("/api/categories", ch.List)
	r.POST("/api/publish", pub.Publish)

	r.NoRoute(spaHandler("web"))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}

// buildFrontend runs `npm run build` inside ./frontend if web/index.html is absent.
func buildFrontend() {
	if _, err := os.Stat("web/index.html"); err == nil {
		return
	}
	log.Println("frontend: building...")
	cmd := exec.Command("npm", "run", "build")
	cmd.Dir = "frontend"
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatalf("frontend build failed: %v", err)
	}
}

func spaHandler(dir string) gin.HandlerFunc {
	fs := http.Dir(dir)
	fileServer := http.FileServer(fs)
	return func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/api/") {
			c.Status(http.StatusNotFound)
			return
		}
		f, err := fs.Open(c.Request.URL.Path)
		if err == nil {
			f.Close()
			fileServer.ServeHTTP(c.Writer, c.Request)
			return
		}
		index, err := os.ReadFile(dir + "/index.html")
		if err != nil {
			c.String(http.StatusServiceUnavailable, "frontend not built")
			return
		}
		c.Data(http.StatusOK, "text/html; charset=utf-8", index)
	}
}
