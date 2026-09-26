// Package employees is the module that owns employees, their accounts and
// sessions. It is the identity provider of the monolith: it issues the access
// tokens and offers the other modules a contracts.TokenVerifier to check them.
//
// Everything but this file and the contracts package lives under internal/, so
// the Go compiler itself keeps the other modules away from its domain,
// use cases and tables.
package employees

import (
	"context"
	"crypto/ed25519"
	"fmt"
	"io/fs"
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/franciscoHonorat/Sys-Called/backend/internal/platform/eventbus"
	"github.com/franciscoHonorat/Sys-Called/backend/modules/employees/contracts"
	httpapi "github.com/franciscoHonorat/Sys-Called/backend/modules/employees/internal/adapters/in/http"
	"github.com/franciscoHonorat/Sys-Called/backend/modules/employees/internal/adapters/in/outbox"
	busadapter "github.com/franciscoHonorat/Sys-Called/backend/modules/employees/internal/adapters/out/eventbus"
	"github.com/franciscoHonorat/Sys-Called/backend/modules/employees/internal/adapters/out/postgres"
	"github.com/franciscoHonorat/Sys-Called/backend/modules/employees/internal/adapters/out/security"
	"github.com/franciscoHonorat/Sys-Called/backend/modules/employees/internal/application"
)

// Schema is the Postgres schema the module owns.
const Schema = "employees"

const (
	outboxRelayInterval = time.Second
	accessTokenTTL      = 15 * time.Minute
	refreshTokenTTL     = 7 * 24 * time.Hour
	demoPassword        = "senha123"
)

var demoEmployees = []application.RegisterEmployeeInput{
	{ID: "agent-1", Name: "Ana Souza", Username: "ana", Password: demoPassword, Role: "support"},
	{ID: "agent-2", Name: "Bruno Lima", Username: "bruno", Password: demoPassword, Role: "support"},
	{ID: "agent-3", Name: "Carla Melo", Username: "carla", Password: demoPassword, Role: "support"},
	{ID: "admin-1", Name: "Administradora", Username: "admin", Password: demoPassword, Role: "admin"},
	{ID: "user-1", Name: "Usuário Padrão", Username: "usuario", Password: demoPassword, Role: "user"},
}

// Migrations returns the module's SQL migrations, applied to Schema.
func Migrations() fs.FS {
	sub, err := fs.Sub(postgres.Migrations, "migrations")
	if err != nil {
		panic(err)
	}
	return sub
}

type Config struct {
	Pool *pgxpool.Pool
	Bus  *eventbus.Bus
	// JWTPrivateKey is the Ed25519 key (PKCS8 PEM) that signs access tokens.
	// Empty means an ephemeral key, which RequireSigningKey forbids.
	JWTPrivateKey     string
	RequireSigningKey bool
	// SeedDemoUsers creates the development users (password "senha123").
	SeedDemoUsers bool
	// CookiePath is where the auth routes are mounted, e.g. /api/employees/auth.
	CookiePath string
	Logger     *slog.Logger
}

type Module struct {
	handler      *httpapi.Handler
	relay        *outbox.Relay
	authenticate *application.AuthenticateUseCase
}

var _ contracts.TokenVerifier = (*Module)(nil)

func New(ctx context.Context, cfg Config) (*Module, error) {
	signingKey, generated, err := security.LoadSigningKey(cfg.JWTPrivateKey)
	if err != nil {
		return nil, fmt.Errorf("employees: loading JWT signing key: %w", err)
	}
	if generated {
		if cfg.RequireSigningKey {
			return nil, fmt.Errorf("employees: JWT_PRIVATE_KEY is required in production")
		}
		cfg.Logger.Warn("JWT_PRIVATE_KEY not set: using an ephemeral signing key, sessions will not survive a restart")
	}

	repo := postgres.NewEmployeeRepository(cfg.Pool)
	hasher := security.NewBcryptHasher()

	if cfg.SeedDemoUsers {
		register := application.NewRegisterEmployeeUseCase(repo, hasher)
		for _, seed := range demoEmployees {
			if err := register.Execute(ctx, seed); err != nil {
				return nil, fmt.Errorf("employees: seeding %s: %w", seed.ID, err)
			}
		}
		cfg.Logger.Info("demo users seeded", "count", len(demoEmployees))
	}

	tokenIssuer := security.NewJWTIssuer(signingKey, accessTokenTTL)
	sessions := application.NewSessions(
		tokenIssuer,
		security.NewRefreshTokenGenerator(),
		postgres.NewRefreshTokenStore(cfg.Pool),
		refreshTokenTTL,
		time.Now,
	)
	authenticate := application.NewAuthenticateUseCase(security.NewJWTVerifier(signingKey.Public().(ed25519.PublicKey)))

	handler := httpapi.NewHandler(httpapi.UseCases{
		Authenticate:           authenticate,
		ListEmployees:          application.NewListEmployeesUseCase(repo),
		Login:                  application.NewLoginUseCase(repo, hasher, sessions),
		RefreshSession:         application.NewRefreshSessionUseCase(repo, sessions),
		Logout:                 application.NewLogoutUseCase(sessions),
		PublicKeys:             application.NewGetPublicKeysUseCase(tokenIssuer),
		SignUp:                 application.NewSignUpUseCase(repo, hasher),
		RequestPasswordReset:   application.NewRequestPasswordResetUseCase(repo),
		ApproveEmployee:        application.NewApproveEmployeeUseCase(repo),
		IssueTemporaryPassword: application.NewIssueTemporaryPasswordUseCase(repo, hasher, security.NewTemporaryPasswordGenerator(), sessions),
		ChangePassword:         application.NewChangePasswordUseCase(repo, hasher, sessions),
	}, httpapi.WithCookiePath(cfg.CookiePath))

	relay := outbox.NewRelay(application.NewPublishPendingEventsUseCase(
		postgres.NewOutboxStore(cfg.Pool),
		busadapter.NewPublisher(cfg.Bus),
	))

	return &Module{handler: handler, relay: relay, authenticate: authenticate}, nil
}

// RegisterRoutes mounts the module's HTTP API on r.
func (m *Module) RegisterRoutes(r gin.IRouter) {
	httpapi.RegisterRoutes(r, m.handler)
}

// Run drives the module's background work, the outbox relay, until ctx ends.
func (m *Module) Run(ctx context.Context) {
	m.relay.Run(ctx, outboxRelayInterval)
}

// VerifyAccessToken implements contracts.TokenVerifier for the other modules.
func (m *Module) VerifyAccessToken(_ context.Context, token string) (contracts.Principal, error) {
	caller, err := m.authenticate.Execute(token)
	if err != nil {
		return contracts.Principal{}, err
	}
	return contracts.Principal{ID: caller.ID, Name: caller.Name, Role: caller.Role.String(), MustChangePassword: caller.MustChangePassword}, nil
}
