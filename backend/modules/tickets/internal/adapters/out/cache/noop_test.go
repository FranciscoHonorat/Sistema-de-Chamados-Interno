package cache_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/franciscoHonorat/Sys-Called/backend/modules/tickets/internal/adapters/out/cache"
	"github.com/franciscoHonorat/Sys-Called/backend/modules/tickets/internal/domain/ticket"
)

func TestNoopTicketCache(t *testing.T) {
	c := cache.NoopTicketCache{}
	c.Set(context.Background(), &ticket.Ticket{})

	_, ok := c.Get(context.Background(), uuid.New())

	assert.False(t, ok)
}
