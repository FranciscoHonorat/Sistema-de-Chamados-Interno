package httpapi

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// NewRouter serves the module on its own, the way the tests exercise it.
func NewRouter(h *Handler) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())

	r.GET("/healthz", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
	RegisterRoutes(r, h)

	return r
}

// RegisterRoutes mounts the module's routes on r, which in the monolith is the
// /api/employees group of the shared server.
func RegisterRoutes(r gin.IRouter, h *Handler) {
	authenticated := RequireAuthentication(h.useCases.Authenticate)
	changingPassword := RequireAuthenticationToChangePassword(h.useCases.Authenticate)

	r.GET("/employees", authenticated, h.ListEmployees)
	r.POST("/employees/:id/approve", authenticated, h.ApproveEmployee)
	r.POST("/employees/:id/temporary-password", authenticated, h.IssueTemporaryPassword)
	r.POST("/auth/signup", h.SignUp)
	r.POST("/auth/password-reset-requests", h.RequestPasswordReset)
	r.POST("/auth/change-password", changingPassword, h.ChangePassword)
	r.POST("/auth/login", h.Login)
	r.POST("/auth/refresh", h.RefreshSession)
	r.POST("/auth/logout", h.Logout)
	r.GET("/.well-known/jwks.json", h.JWKS)
}
