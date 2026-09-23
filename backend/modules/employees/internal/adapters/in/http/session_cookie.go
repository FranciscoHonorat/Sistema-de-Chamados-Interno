package httpapi

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	refreshCookieName = "refresh_token"
	defaultCookiePath = "/auth"
)

func (h *Handler) setRefreshCookie(c *gin.Context, value string, maxAge int) {
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     refreshCookieName,
		Value:    value,
		Path:     h.cookiePath,
		MaxAge:   maxAge,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})
}

func (h *Handler) clearRefreshCookie(c *gin.Context) {
	h.setRefreshCookie(c, "", -1)
}
