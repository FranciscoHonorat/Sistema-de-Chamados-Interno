// Package contracts is the published language of the employees module: the
// only package other modules may import from it. It holds the integration
// events the module emits through the outbox and the port it offers to
// authenticate the bearer of an access token. It depends on nothing but the
// standard library, so importing it never drags the module's internals along.
package contracts

import "context"

// Event types published on the event bus.
const (
	EventEmployeeRegistered     = "EmployeeRegistered"
	EventEmployeeSignedUp       = "EmployeeSignedUp"
	EventPasswordResetRequested = "PasswordResetRequested"
)

// EmployeeRegistered is published when an employee is created by the system
// (seed or administrator). Role is one of "user", "support" and "admin".
type EmployeeRegistered struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Role string `json:"role"`
}

// EmployeeSignedUp is published when someone creates their own account, which
// then waits for an administrator's approval.
type EmployeeSignedUp struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// PasswordResetRequested is published when an employee asks the administrator
// for a temporary password.
type PasswordResetRequested struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Principal is the authenticated bearer of an access token.
type Principal struct {
	ID   string
	Name string
	Role string

	MustChangePassword bool
}

// TokenVerifier validates an access token issued by the employees module.
type TokenVerifier interface {
	VerifyAccessToken(ctx context.Context, token string) (Principal, error)
}
