package valueobjects_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	status "github.com/franciscoHonorat/Sys-Called/backend/modules/tickets/internal/domain/valueobjects"
)

func TestStatus(t *testing.T) {
	t.Run("should create a new Status with a valid string", func(t *testing.T) {
		statusStr := "Open"
		statusObj := status.Status(statusStr)

		assert.Equal(t, statusStr, statusObj.String())
	})

	t.Run("should check validity of Status", func(t *testing.T) {
		validStatus := status.TicketStatusOpen
		invalidStatus := status.Status("Invalid")

		assert.True(t, validStatus.IsValid())
		assert.False(t, invalidStatus.IsValid())
	})

}
