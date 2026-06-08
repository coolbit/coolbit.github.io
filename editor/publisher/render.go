package publisher

import (
	"blog/store"
	"bytes"
	"fmt"
	"html/template"
)

type IndexData struct {
	Posts      []store.Post
	Page       int
	TotalPages int
	PrevURL    string
	NextURL    string
	PostsBase  string
	HomeURL    string
}

var indexTpl = template.Must(template.New("index").Parse(`<!DOCTYPE html>
<html lang="zh">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>{{if gt .Page 1}}Page {{.Page}} — {{end}}Dawn's Blog</title>
<link rel="stylesheet" href="{{.HomeURL}}/style.css">
</head>
<body>
<header class="site-header">
  <div class="container hd-inner">
    <a href="{{.HomeURL}}/" class="brand">Dawn's Blog</a>
  </div>
</header>
<main class="container">
{{if not .Posts}}<p class="empty">No posts yet.</p>{{end}}
<div class="post-list">
{{range .Posts}}
<article class="post-card">
  <div class="post-meta">
    {{if .Category}}<span class="cat">{{.Category}}</span>{{end}}
    <span class="date">{{.CreatedAt.Format "January 2, 2006"}}</span>
  </div>
  <h2 class="post-title"><a href="{{$.PostsBase}}{{.ID}}.html">{{.Title}}</a></h2>
  {{if .Summary}}<p class="summary">{{.Summary}}</p>{{end}}
</article>
{{end}}
</div>
{{if gt .TotalPages 1}}
<nav class="pagination">
  {{if .PrevURL}}
    <a href="{{.PrevURL}}" class="page-btn">← Newer</a>
  {{else}}
    <span class="page-btn disabled">← Newer</span>
  {{end}}
  <span class="page-info">{{.Page}} / {{.TotalPages}}</span>
  {{if .NextURL}}
    <a href="{{.NextURL}}" class="page-btn">Older →</a>
  {{else}}
    <span class="page-btn disabled">Older →</span>
  {{end}}
</nav>
{{end}}
</main>
</body>
</html>`))

func RenderIndex(data IndexData) (string, error) {
	var buf bytes.Buffer
	err := indexTpl.Execute(&buf, data)
	return buf.String(), err
}

func PaginateIndex(posts []store.Post, pageSize int) []struct {
	Path string
	Data IndexData
} {
	total := len(posts)
	totalPages := (total + pageSize - 1) / pageSize
	if totalPages == 0 {
		totalPages = 1
	}

	var pages []struct {
		Path string
		Data IndexData
	}

	for p := 1; p <= totalPages; p++ {
		start := (p - 1) * pageSize
		end := start + pageSize
		if end > total {
			end = total
		}

		var prevURL, nextURL, postsBase, homeURL string

		if p == 1 {
			postsBase = "posts/"
			homeURL = "."
			if totalPages > 1 {
				nextURL = "page/2.html"
			}
		} else {
			postsBase = "../posts/"
			homeURL = ".."
			if p == 2 {
				prevURL = "../index.html"
			} else {
				prevURL = fmt.Sprintf("%d.html", p-1)
			}
			if p < totalPages {
				nextURL = fmt.Sprintf("%d.html", p+1)
			}
		}

		path := "index.html"
		if p > 1 {
			path = fmt.Sprintf("page/%d.html", p)
		}

		pages = append(pages, struct {
			Path string
			Data IndexData
		}{
			Path: path,
			Data: IndexData{
				Posts:      posts[start:end],
				Page:       p,
				TotalPages: totalPages,
				PrevURL:    prevURL,
				NextURL:    nextURL,
				PostsBase:  postsBase,
				HomeURL:    homeURL,
			},
		})
	}

	return pages
}

const StyleCSS = `*,*::before,*::after{box-sizing:border-box;margin:0;padding:0}
body{font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',sans-serif;color:#111;background:#fff;line-height:1.6;font-size:16px}
a{color:inherit;text-decoration:none}
a:hover{color:#0070f3}
.container{max-width:760px;margin:0 auto;padding:0 2rem}
.site-header{border-bottom:1px solid #ebebeb;padding:1rem 0}
.hd-inner{display:flex;align-items:center;justify-content:space-between}
.brand{font-size:1.125rem;font-weight:800;letter-spacing:-.02em}
.back{font-size:.875rem;color:#999}
.back:hover{color:#111}
main{padding:3rem 0}
.empty{color:#bbb;text-align:center;padding:3rem 0}
.post-list{display:flex;flex-direction:column}
.post-card{padding:2rem 0;border-bottom:1px solid #f3f3f3}
.post-card:last-child{border-bottom:none}
.post-meta{display:flex;align-items:center;gap:.75rem;margin-bottom:.625rem}
.cat{font-size:.75rem;font-weight:600;text-transform:uppercase;letter-spacing:.05em;color:#0070f3}
.date{font-size:.8125rem;color:#bbb}
.post-title{font-size:1.5rem;font-weight:700;letter-spacing:-.02em;line-height:1.3;margin-bottom:.625rem}
.post-title a:hover{color:#0070f3}
.summary{font-size:.9375rem;color:#555;line-height:1.7;display:-webkit-box;-webkit-line-clamp:2;-webkit-box-orient:vertical;overflow:hidden}
.post-detail{padding:2rem 0}
.detail-title{font-size:2.5rem;font-weight:800;letter-spacing:-.04em;line-height:1.15;margin-bottom:1rem}
.lead{font-size:1.125rem;color:#555;line-height:1.75;margin-bottom:1.5rem}
.divider{border:none;border-top:1px solid #f0f0f0;margin:2rem 0}
.md-body{font-size:1rem;line-height:1.8}
.md-body h1,.md-body h2,.md-body h3,.md-body h4{font-weight:700;letter-spacing:-.02em;margin:1.5em 0 .75em;line-height:1.3}
.md-body h1{font-size:1.875rem}.md-body h2{font-size:1.5rem}.md-body h3{font-size:1.25rem}
.md-body p{margin-bottom:1.25rem}
.md-body a{color:#0070f3;text-decoration:underline}
.md-body ul,.md-body ol{margin:1rem 0 1rem 1.5rem}
.md-body li{margin-bottom:.5rem}
.md-body code{font-family:'SFMono-Regular',Consolas,monospace;font-size:.875em;background:#f5f5f5;padding:.125em .375em;border-radius:4px}
.md-body pre{background:#1a1a2e;color:#e0e0e0;padding:1.25rem;border-radius:8px;overflow-x:auto;margin:1.5rem 0}
.md-body pre code{background:none;padding:0;color:inherit;font-size:.875rem}
.md-body blockquote{border-left:3px solid #e5e5e5;padding:.5rem 1rem;color:#666;margin:1.5rem 0}
.md-body table{width:100%;border-collapse:collapse;margin:1.5rem 0}
.md-body th,.md-body td{padding:.5rem 1rem;border:1px solid #e5e5e5;text-align:left}
.md-body th{background:#f9f9f9;font-weight:600}
.md-body img{max-width:100%;border-radius:8px}
.md-body hr{border:none;border-top:1px solid #f0f0f0;margin:2rem 0}
.pagination{display:flex;align-items:center;justify-content:center;gap:1.25rem;margin-top:3rem;padding-top:2rem;border-top:1px solid #f0f0f0}
.page-btn{padding:.4rem 1rem;border:1px solid #e5e5e5;border-radius:6px;font-size:.875rem;color:#555;transition:all .15s;text-decoration:none!important}
.page-btn:hover{border-color:#111;color:#111}
.page-btn.disabled{opacity:.35;pointer-events:none;cursor:default}
.page-info{font-size:.875rem;color:#bbb;min-width:60px;text-align:center}
`
