package valueobjects_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	domainErr "github.com/franciscoHonorat/Sys-Called/backend/modules/tickets/internal/domain/domain-errors"
	authorID "github.com/franciscoHonorat/Sys-Called/backend/modules/tickets/internal/domain/valueobjects"
)

func TestAuthorID(t *testing.T) {
	t.Run("should create a new AuthorID with a valid string", func(t *testing.T) {
		idStr := "valid-author-id"
		idObj, err := authorID.NewAuthorID(idStr)

		assert.NoError(t, err)
		assert.Equal(t, idStr, idObj.GetAuthorID())
	})

	t.Run("should return an error when creating an AuthorID with an empty string", func(t *testing.T) {
		_, err := authorID.NewAuthorID("")

		assert.Error(t, err)
		assert.ErrorIs(t, err, domainErr.ErrInvalidAuthorID)
	})

}
