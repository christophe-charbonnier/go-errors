package code

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewLocation(t *testing.T) {
	t.Run("absolute", func(t *testing.T) {
		location := NewLocation("/Users/CharbonC/dev/maers-brt-bc-atos/v3/internal/errors/errors.go:17")
		assert.Equal(t, "/Users/CharbonC/dev/maers-brt-bc-atos/v3/internal/errors", location.Dir)
		assert.Equal(t, "errors.go", location.File)
		assert.Equal(t, 17, location.Line)
	})

	t.Run("relative", func(t *testing.T) {
		location := NewLocation("errors/errors.go:17")
		assert.Equal(t, "errors", location.Dir)
		assert.Equal(t, "errors.go", location.File)
		assert.Equal(t, 17, location.Line)
	})

	t.Run("no-dir", func(t *testing.T) {
		location := NewLocation("log_test.go:32")
		assert.Equal(t, "", location.Dir)
		assert.Equal(t, "log_test.go", location.File)
		assert.Equal(t, 32, location.Line)
	})
}
