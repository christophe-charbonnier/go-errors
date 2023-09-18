package errors

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/pterm/pterm"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

func init() {
	pterm.DisableColor()
}

func Test_New(t *testing.T) {
	err := New("test error\nsecond line\nthird line")
	assert.Equal(t, "moodys.com/medor/internal/errors", err.Frames()[0].Code.Package)
	assert.Equal(t, "Test_New", err.Frames()[0].Code.Func)
	fmt.Println(err)
}

func Test_WithError(t *testing.T) {
	err := WithError(fmt.Errorf("wrapped error"))
	assert.Equal(t, "moodys.com/medor/internal/errors", err.Frames()[0].Code.Package)
	assert.Equal(t, "Test_WithError", err.Frames()[0].Code.Func)
	fmt.Println(err)
}

func Test_Json(t *testing.T) {
	err := New("foo")
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	enc.Encode(err)
}

func Test_Log(t *testing.T) {
	err := New("foo")
	log.Error().Msgf("error=%v", err)
}
