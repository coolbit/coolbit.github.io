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

	mediaDir := os.Getenv("MEDIA_DIR")
	if mediaDir == "" {
		mediaDir = "../media"
	}

	r := gin.Default()
	r.Use(cors.Default())

	ph := handlers.NewPostHandler(s, mediaDir)
	pub := handlers.NewPublishHandler(s)
	uh := handlers.NewUploadHandler(mediaDir)
	mh := handlers.NewMediaListHandler(mediaDir)

	posts := r.Group("/api/posts")
	{
		posts.GET("", ph.List)
		posts.GET("/:id", ph.Get)
		posts.POST("", ph.Create)
		posts.PUT("/:id", ph.Update)
		posts.DELETE("/:id", ph.Delete)
	}

	r.POST("/api/upload/:id", uh.Upload)
	r.GET("/api/media/:id", mh.List)
	r.POST("/api/publish", pub.Publish)
	r.Static("/media", mediaDir)

	r.NoRoute(spaHandler("web"))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}

// buildFrontend installs npm deps (if needed) and builds the frontend when web/index.html is absent.
func buildFrontend() {
	if _, err := os.Stat("web/index.html"); err == nil {
		return
	}
	if _, err := os.Stat("frontend/node_modules"); err != nil {
		log.Println("frontend: installing dependencies...")
		install := exec.Command("npm", "install")
		install.Dir = "frontend"
		install.Stdout = os.Stdout
		install.Stderr = os.Stderr
		if err := install.Run(); err != nil {
			log.Fatalf("npm install failed: %v", err)
		}
	}
	log.Println("frontend: building...")
	build := exec.Command("npm", "run", "build")
	build.Dir = "frontend"
	build.Stdout = os.Stdout
	build.Stderr = os.Stderr
	if err := build.Run(); err != nil {
		log.Fatalf("npm run build failed: %v", err)
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
