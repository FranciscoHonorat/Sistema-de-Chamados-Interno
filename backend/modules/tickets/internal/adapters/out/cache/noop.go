package cache

import (
	"context"

	"github.com/google/uuid"

	"github.com/franciscoHonorat/Sys-Called/backend/modules/tickets/internal/application/port/out"
	"github.com/franciscoHonorat/Sys-Called/backend/modules/tickets/internal/domain/ticket"
)

var _ out.TicketCache = NoopTicketCache{}

// NoopTicketCache never remembers a ticket, so every read rebuilds it from
// the event store. It is the cache to use with more than one replica, where a
// per-instance memory cache would serve a ticket another replica changed.
type NoopTicketCache struct{}

func (NoopTicketCache) Get(context.Context, uuid.UUID) (*ticket.Ticket, bool) { return nil, false }

func (NoopTicketCache) Set(context.Context, *ticket.Ticket) {}
