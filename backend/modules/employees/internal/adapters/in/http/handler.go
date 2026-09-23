package httpapi

import (
	"github.com/gin-gonic/gin"

	"github.com/franciscoHonorat/Sys-Called/backend/modules/employees/internal/application"
)

type UseCases struct {
	Authenticate           *application.AuthenticateUseCase
	ListEmployees          *application.ListEmployeesUseCase
	Login                  *application.LoginUseCase
	RefreshSession         *application.RefreshSessionUseCase
	Logout                 *application.LogoutUseCase
	PublicKeys             *application.GetPublicKeysUseCase
	SignUp                 *application.SignUpUseCase
	RequestPasswordReset   *application.RequestPasswordResetUseCase
	ApproveEmployee        *application.ApproveEmployeeUseCase
	IssueTemporaryPassword *application.IssueTemporaryPasswordUseCase
	ChangePassword         *application.ChangePasswordUseCase
}

type Handler struct {
	useCases   UseCases
	cookiePath string
}

type Option func(*Handler)

// WithCookiePath restricts the refresh cookie to the path the auth routes are
// mounted on, so the browser only sends it to them.
func WithCookiePath(path string) Option {
	return func(h *Handler) { h.cookiePath = path }
}

func NewHandler(useCases UseCases, options ...Option) *Handler {
	h := &Handler{useCases: useCases, cookiePath: defaultCookiePath}
	for _, option := range options {
		option(h)
	}
	return h
}

func (h *Handler) ListEmployees(c *gin.Context) {
	output, err := h.useCases.ListEmployees.Execute(c.Request.Context(), CallerFrom(c))
	respondJSON(c, output, err)
}
