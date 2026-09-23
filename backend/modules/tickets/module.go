// Package tickets is the module that owns support tickets, their event-sourced
// history, the list of support agents and the notifications. It learns about
// employees only through the events and the token verifier the employees
// module publishes in its contracts package.
package tickets

import (
	"io/fs"
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/franciscoHonorat/Sys-Called/backend/internal/platform/eventbus"
	"github.com/franciscoHonorat/Sys-Called/backend/modules/employees/contracts"
	"github.com/franciscoHonorat/Sys-Called/backend/modules/tickets/internal/adapters/in/events"
	httpapi "github.com/franciscoHonorat/Sys-Called/backend/modules/tickets/internal/adapters/in/http"
	"github.com/franciscoHonorat/Sys-Called/backend/modules/tickets/internal/adapters/out/cache"
	"github.com/franciscoHonorat/Sys-Called/backend/modules/tickets/internal/adapters/out/employees"
	"github.com/franciscoHonorat/Sys-Called/backend/modules/tickets/internal/adapters/out/postgres"
	"github.com/franciscoHonorat/Sys-Called/backend/modules/tickets/internal/application"
	"github.com/franciscoHonorat/Sys-Called/backend/modules/tickets/internal/application/auth"
	"github.com/franciscoHonorat/Sys-Called/backend/modules/tickets/internal/application/command"
	"github.com/franciscoHonorat/Sys-Called/backend/modules/tickets/internal/application/port/out"
	"github.com/franciscoHonorat/Sys-Called/backend/modules/tickets/internal/application/query"
)

// Schema is the Postgres schema the module owns.
const Schema = "tickets"

// Migrations returns the module's SQL migrations, applied to Schema.
func Migrations() fs.FS {
	sub, err := fs.Sub(postgres.Migrations, "migrations")
	if err != nil {
		panic(err)
	}
	return sub
}

type Config struct {
	Pool     *pgxpool.Pool
	Bus      *eventbus.Bus
	Identity contracts.TokenVerifier
	// CacheTickets keeps rebuilt tickets in memory. Only safe with a single
	// replica: the cache is per instance and is not invalidated by the others.
	CacheTickets bool
	Logger       *slog.Logger
}

type Module struct {
	handler *httpapi.Handler
}

// New wires the module and subscribes it to the employees module's events.
func New(cfg Config) *Module {
	notifications := postgres.NewNotificationStore(cfg.Pool)
	store := application.NewNotifyingEventStore(postgres.NewEventStore(cfg.Pool), notifications)
	var ticketCache out.TicketCache = cache.NoopTicketCache{}
	if cfg.CacheTickets {
		ticketCache = cache.NewInMemoryTicketCache()
	}
	responsibles := postgres.NewResponsibleDirectory(cfg.Pool)

	events.NewEmployeesSubscriber(
		command.NewSyncResponsibleUseCase(responsibles),
		command.NewNotifyAccountRequestUseCase(notifications, time.Now),
	).Subscribe(cfg.Bus)

	handler := httpapi.NewHandler(httpapi.UseCases{
		Authenticate:           auth.NewAuthenticateUseCase(employees.NewTokenVerifier(cfg.Identity)),
		OpenTicket:             command.NewOpenTicketUseCase(store, ticketCache, responsibles),
		GetTicket:              query.NewGetTicketUseCase(store, ticketCache),
		ListTickets:            query.NewListTicketsUseCase(store, ticketCache),
		EditTicket:             command.NewEditTicketUseCase(store, ticketCache),
		AssignTicket:           command.NewAssignTicketUseCase(store, ticketCache, responsibles),
		AutoAssignTicket:       command.NewAutoAssignTicketUseCase(store, ticketCache, responsibles),
		ChangeTicketPriority:   command.NewChangeTicketPriorityUseCase(store, ticketCache),
		MoveTicketToInProgress: command.NewMoveTicketToInProgressUseCase(store, ticketCache),
		CloseTicket:            command.NewCloseTicketUseCase(store, ticketCache),
		AddTicketResponse:      command.NewAddTicketResponseUseCase(store, ticketCache),
		ListResponsibles:       query.NewListResponsiblesUseCase(responsibles),
		ListNotifications:      query.NewListNotificationsUseCase(notifications),
		MarkNotificationsRead:  command.NewMarkNotificationsReadUseCase(notifications, time.Now),
		SupportWorkload:        query.NewSupportWorkloadUseCase(store, ticketCache, responsibles),
	})

	return &Module{handler: handler}
}

// RegisterRoutes mounts the module's HTTP API on r.
func (m *Module) RegisterRoutes(r gin.IRouter) {
	httpapi.RegisterRoutes(r, m.handler)
}
