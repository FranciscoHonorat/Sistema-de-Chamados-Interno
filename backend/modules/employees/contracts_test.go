package employees_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/franciscoHonorat/Sys-Called/backend/modules/employees/contracts"
	"github.com/franciscoHonorat/Sys-Called/backend/modules/employees/internal/domain/employee"
)

// The outbox stores the domain events as they are, and the other modules
// decode them with the published contracts. This test is what keeps the two
// from drifting apart.
func TestDomainEventsMatchThePublishedContracts(t *testing.T) {
	t.Run("EmployeeRegistered", func(t *testing.T) {
		event := employee.Registered{ID: "agent-1", Name: "Ana Souza", Role: "support"}
		assert.Equal(t, contracts.EventEmployeeRegistered, event.EventType())
		assert.Equal(t, contracts.EmployeeRegistered{ID: "agent-1", Name: "Ana Souza", Role: "support"}, roundTrip[contracts.EmployeeRegistered](t, event))
	})

	t.Run("EmployeeSignedUp", func(t *testing.T) {
		event := employee.SignedUp{ID: "user-9", Name: "Maria Lima"}
		assert.Equal(t, contracts.EventEmployeeSignedUp, event.EventType())
		assert.Equal(t, contracts.EmployeeSignedUp{ID: "user-9", Name: "Maria Lima"}, roundTrip[contracts.EmployeeSignedUp](t, event))
	})

	t.Run("PasswordResetRequested", func(t *testing.T) {
		event := employee.PasswordResetRequested{ID: "user-1", Name: "Usuário Padrão"}
		assert.Equal(t, contracts.EventPasswordResetRequested, event.EventType())
		assert.Equal(t, contracts.PasswordResetRequested{ID: "user-1", Name: "Usuário Padrão"}, roundTrip[contracts.PasswordResetRequested](t, event))
	})
}

func roundTrip[T any](t *testing.T, event employee.Event) T {
	t.Helper()
	payload, err := json.Marshal(event)
	require.NoError(t, err)
	var decoded T
	require.NoError(t, json.Unmarshal(payload, &decoded))
	return decoded
}
