package out

import (
	"github.com/franciscoHonorat/Sys-Called/backend/modules/employees/internal/domain/employee"
)

type Caller struct {
	ID   string
	Name string
	Role employee.Role

	MustChangePassword bool
}

type TokenVerifier interface {
	Verify(token string) (Caller, error)
}
