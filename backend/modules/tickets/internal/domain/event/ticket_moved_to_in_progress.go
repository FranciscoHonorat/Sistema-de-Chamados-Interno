package event

import (
	"github.com/franciscoHonorat/Sys-Called/backend/modules/tickets/internal/domain/valueobjects"
)

type TicketMovedToInProgress struct {
	baseEvent
}

func NewTicketMovedToInProgress(id *valueobjects.ID) TicketMovedToInProgress {
	return TicketMovedToInProgress{
		baseEvent: newBaseEvent(id.GetID()),
	}
}

func (TicketMovedToInProgress) EventName() string {
	return "TicketMovedToInProgress"
}

func init() {
	registerEvent[TicketMovedToInProgress]()
}
