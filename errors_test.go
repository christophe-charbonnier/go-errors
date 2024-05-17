package errors

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_New(t *testing.T) {
	err := New("test error\nsecond line\nthird line")
	assert.Equal(t, "test error\nsecond line\nthird line", err.Message)
	assert.Equal(t, "github.com/moodysanalytics/maers-brt-bc-atos/internal/errors", err.Frames[0].Code.Package)
	assert.Equal(t, "Test_New", err.Frames[0].Code.Func)
	fmt.Println(err)
}

func Test_WithError(t *testing.T) {
	t.Run("", func(t *testing.T) {
		err := WithError(fmt.Errorf("wrapped error"))
		assert.Equal(t, "github.com/moodysanalytics/maers-brt-bc-atos/internal/errors", err.Frames[0].Code.Package)
		assert.Equal(t, "Test_WithError.func1", err.Frames[0].Code.Func)
		fmt.Println(err)
	})

	t.Run("", func(t *testing.T) {
		err := WithError(fmt.Errorf("wrapped error"))
		err = WithError(err)
		fmt.Println(err)
	})
}

func Test_WithRecover(t *testing.T) {
	t.Run("panic", func(t *testing.T) {
		defer func() { fmt.Println(WithRecover(recover())) }()
		panic("fake runtime error")
	})

	t.Run("nil-dereference", func(t *testing.T) {
		defer func() { fmt.Println(WithRecover(recover())) }()
		*(*int)(nil) = 1
	})

	t.Run("division-by-0", func(t *testing.T) {
		defer func() { fmt.Println(WithRecover(recover())) }()
		var d int
		_ = 1 / d
	})
}

type Test struct {
	Value int
}

func (t *Test) PointerReceiver() error {
	return New("in PointerReceiver")
}

func (t Test) Receiver() error {
	return New("in Receiver")
}

func (t *Test) Panic() (err error) {
	defer func() {
		err = WithRecover(recover())
	}()

	panic("fake runtime error")
}

func (t *Test) NilDereference() (err error) {
	defer func() {
		err = WithRecover(recover())
	}()

	*(*int)(nil) = 1
	return
}

func Test_Json(t *testing.T) {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	assert.NoError(t, enc.Encode(New("foo")))
}

func Test_Log(t *testing.T) {
	// err := New("foo")
}
