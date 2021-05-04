package model

import (
	"testing"
	"time"

	"github.com/adhocore/urlsh/common"
)

func TestURL_IsActive(t *testing.T) {
	t.Run("is active - deleted", func(t *testing.T) {
		model := URL{Deleted: true}
		if model.IsActive() {
			t.Errorf("should not be active if deleted")
		}
	})

	t.Run("is active - expired", func(t *testing.T) {
		past, _ := time.ParseInLocation(common.DateLayout, "2000-01-01 00:00:00", time.UTC)
		model := URL{Deleted: false, ExpiresOn: past}
		if model.IsActive() {
			t.Errorf("should not be active if expired")
		}
	})

	t.Run("is active - OK", func(t *testing.T) {
		future, _ := time.ParseInLocation(common.DateLayout, "3000-01-01 00:00:00", time.UTC)
		model := URL{Deleted: false, ExpiresOn: future}
		if !model.IsActive() {
			t.Errorf("should be active if not deleted or expired")
		}
	})
}
