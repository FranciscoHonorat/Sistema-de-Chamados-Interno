package httpserver

import (
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

// serveSPA serves the built frontend from dir: files that exist are served as
// they are (the hashed ones under /assets/ cached for a year), and every other
// path gets index.html, so the Vue router can handle deep links. Unknown /api
// paths still answer with a JSON 404.
func serveSPA(dir string) gin.HandlerFunc {
	index := filepath.Join(dir, "index.html")
	files := http.Dir(dir)

	return func(c *gin.Context) {
		p := c.Request.URL.Path
		if strings.HasPrefix(p, "/api/") || p == "/api" {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		if c.Request.Method != http.MethodGet && c.Request.Method != http.MethodHead {
			c.Status(http.StatusMethodNotAllowed)
			return
		}

		clean := path.Clean("/" + p)
		if clean != "/" && fileExists(files, clean) {
			if strings.HasPrefix(clean, "/assets/") {
				c.Header("Cache-Control", "public, max-age=31536000, immutable")
			} else {
				c.Header("Cache-Control", "no-cache")
			}
			c.FileFromFS(clean, files)
			return
		}

		c.Header("Cache-Control", "no-cache")
		c.File(index)
	}
}

func fileExists(files http.FileSystem, name string) bool {
	f, err := files.Open(name)
	if err != nil {
		return false
	}
	defer f.Close()
	info, err := f.Stat()
	return err == nil && !info.IsDir()
}

func staticDirUsable(dir string) bool {
	if dir == "" {
		return false
	}
	_, err := os.Stat(filepath.Join(dir, "index.html"))
	return err == nil
}
