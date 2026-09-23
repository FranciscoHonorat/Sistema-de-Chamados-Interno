package httpserver

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Check reports whether a dependency the app needs to serve traffic is usable.
type Check func(ctx context.Context) error

func registerHealth(r *gin.Engine, checks map[string]Check) {
	// Liveness: the process is up and serving HTTP. It never looks at the
	// database, so a database outage does not make Kubernetes restart the pods.
	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Readiness: every dependency answers, so the instance may get traffic.
	r.GET("/readyz", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
		defer cancel()

		status, results := http.StatusOK, gin.H{}
		for name, check := range checks {
			if err := check(ctx); err != nil {
				status = http.StatusServiceUnavailable
				results[name] = err.Error()
				continue
			}
			results[name] = "ok"
		}
		c.JSON(status, gin.H{"status": http.StatusText(status), "checks": results})
	})
}
