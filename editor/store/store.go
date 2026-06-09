package store

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	htmlr "github.com/yuin/goldmark/renderer/html"
)

var gm = goldmark.New(
	goldmark.WithExtensions(
		extension.GFM,
		highlighting.NewHighlighting(
			highlighting.WithStyle("monokai"),
		),
	),
	goldmark.WithRendererOptions(htmlr.WithUnsafe()),
)

func renderMarkdown(src string) template.HTML {
	var buf bytes.Buffer
	gm.Convert([]byte(src), &buf)
	return template.HTML(buf.String())
}

type Post struct {
	ID         uint      `json:"id"`
	Title      string    `json:"title"`
	Summary    string    `json:"summary"`
	CoverImage string    `json:"cover_image,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	Content    string    `json:"content"`
}

type Store struct {
	Dir string
}

func New(dir string) (*Store, error) {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}
	return &Store{Dir: dir}, nil
}

func (s *Store) path(id uint) string {
	return filepath.Join(s.Dir, fmt.Sprintf("%d.html", id))
}

// postTpl renders the published HTML. Post metadata is embedded in a JSON
// script tag so the CMS can read it back without a separate source file.
var postTpl = template.Must(template.New("post").Parse(`<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>{{.Post.Title}} — Dawn's Blog</title>
<link rel="stylesheet" href="../style.css">
<script type="application/json" id="post-data">{{.Data}}</script>
</head>
<body>
<header class="site-header">
  <div class="container hd-inner">
    <a href="../" class="brand">Dawn's Blog</a>
    <a href="../" class="back">← All posts</a>
  </div>
</header>
<main class="container">
<article class="post-detail">
  {{if .Post.CoverImage}}<img src="{{.Post.CoverImage}}" class="cover-img" alt="">{{end}}
  <div class="post-meta">
    <span class="date">{{.Post.CreatedAt.Format "January 2, 2006"}}</span>
  </div>
  <h1 class="detail-title">{{.Post.Title}}</h1>
  {{if .Post.Summary}}<p class="lead">{{.Post.Summary}}</p>{{end}}
  <hr class="divider">
  <div class="md-body">{{.Content}}</div>
</article>
</main>
</body>
</html>`))

type postPage struct {
	Post    Post
	Data    template.JS
	Content template.HTML
}

const scriptOpen = `<script type="application/json" id="post-data">`
const scriptClose = `</script>`

func parse(data []byte) (*Post, error) {
	s := string(data)
	si := strings.Index(s, scriptOpen)
	if si < 0 {
		return nil, fmt.Errorf("post-data not found")
	}
	si += len(scriptOpen)
	ei := strings.Index(s[si:], scriptClose)
	if ei < 0 {
		return nil, fmt.Errorf("unclosed script tag")
	}
	var p Post
	if err := json.Unmarshal([]byte(strings.TrimSpace(s[si:si+ei])), &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func (s *Store) write(p *Post) error {
	// json.Marshal escapes <, >, & as \uXXXX — safe inside <script>.
	jsonData, err := json.Marshal(p)
	if err != nil {
		return err
	}
	var buf bytes.Buffer
	if err := postTpl.Execute(&buf, postPage{
		Post:    *p,
		Data:    template.JS(jsonData),
		Content: renderMarkdown(p.Content),
	}); err != nil {
		return err
	}
	return os.WriteFile(s.path(p.ID), buf.Bytes(), 0644)
}

func (s *Store) nextID() uint {
	entries, _ := os.ReadDir(s.Dir)
	var max uint
	for _, e := range entries {
		n, err := strconv.ParseUint(strings.TrimSuffix(e.Name(), ".html"), 10, 64)
		if err == nil && uint(n) > max {
			max = uint(n)
		}
	}
	return max + 1
}

func (s *Store) All() ([]Post, error) {
	entries, err := os.ReadDir(s.Dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var posts []Post
	for _, e := range entries {
		if !strings.HasSuffix(e.Name(), ".html") {
			continue
		}
		if _, err := strconv.ParseUint(strings.TrimSuffix(e.Name(), ".html"), 10, 64); err != nil {
			continue
		}
		data, err := os.ReadFile(filepath.Join(s.Dir, e.Name()))
		if err != nil {
			continue
		}
		p, err := parse(data)
		if err != nil {
			continue
		}
		posts = append(posts, *p)
	}
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].CreatedAt.After(posts[j].CreatedAt)
	})
	return posts, nil
}

func (s *Store) Get(id uint) (*Post, error) {
	data, err := os.ReadFile(s.path(id))
	if err != nil {
		return nil, err
	}
	return parse(data)
}

func (s *Store) Create(title, summary, content, coverImage string) (*Post, error) {
	now := time.Now()
	p := &Post{
		ID:         s.nextID(),
		Title:      title,
		Summary:    summary,
		Content:    content,
		CoverImage: coverImage,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	return p, s.write(p)
}

func (s *Store) Update(id uint, title, summary, content, coverImage string) (*Post, error) {
	p, err := s.Get(id)
	if err != nil {
		return nil, err
	}
	p.Title = title
	p.Summary = summary
	p.Content = content
	p.CoverImage = coverImage
	p.UpdatedAt = time.Now()
	return p, s.write(p)
}

func (s *Store) Delete(id uint) error {
	return os.Remove(s.path(id))
}

// RerenderAll rewrites every post HTML using the current template.
// Called by Publish to keep all posts in sync with template changes.
func (s *Store) RerenderAll() error {
	posts, err := s.All()
	if err != nil {
		return err
	}
	for i := range posts {
		if err := s.write(&posts[i]); err != nil {
			return err
		}
	}
	return nil
}

