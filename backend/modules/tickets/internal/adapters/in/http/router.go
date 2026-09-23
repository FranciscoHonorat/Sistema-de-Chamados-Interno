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
// /api/tickets group of the shared server.
func RegisterRoutes(r gin.IRouter, h *Handler) {
	authenticated := RequireAuthentication(h.useCases.Authenticate)
	r.GET("/responsibles", authenticated, h.ListResponsibles)
	r.GET("/responsibles/workload", authenticated, h.SupportWorkload)
	r.GET("/notifications", authenticated, h.ListNotifications)
	r.POST("/notifications/read", authenticated, h.MarkNotificationsRead)

	tickets := r.Group("/tickets", authenticated)
	tickets.POST("", h.OpenTicket)
	tickets.GET("", h.ListTickets)
	tickets.GET("/:id", h.GetTicket)
	tickets.PUT("/:id", h.EditTicket)
	tickets.POST("/:id/assign", h.AssignTicket)
	tickets.POST("/:id/assign/auto", h.AutoAssignTicket)
	tickets.POST("/:id/priority", h.ChangeTicketPriority)
	tickets.POST("/:id/start", h.MoveTicketToInProgress)
	tickets.POST("/:id/close", h.CloseTicket)
	tickets.POST("/:id/responses", h.AddTicketResponse)
}
