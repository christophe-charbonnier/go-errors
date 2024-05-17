package errors

import (
	"fmt"
	"runtime/debug"
	"strings"

	"github.com/christophe-charbonnier/go-errors/code"
	"github.com/pkg/errors"
)

type Frame struct {
	Location code.Location
	Code     code.Code
}

func NewFrame(s string) (f Frame) {
	lines := strings.Split(s, "\n")
	f.Code = code.NewCode(lines[0])
	if len(lines) > 1 {
		f.Location = code.NewLocation(lines[1])
	}

	return
}

func (f Frame) String(colWidth int) string {
	return fmt.Sprintf("%-*v", colWidth, f.Location.String()) + f.Code.String()
}

func FramesFromError(err error) (frames []Frame) {
	type stackTracer interface {
		StackTrace() errors.StackTrace
	}

	if err, ok := err.(stackTracer); ok {
		for _, f := range err.StackTrace() {
			frames = append(frames, NewFrame(fmt.Sprintf("%+v", f)))
		}
	}

	return
}

func FramesFromDebug() (frames []Frame) {
	lines := strings.Split(string(debug.Stack()), "\n")

	var isAfterPanic bool
	for i := 1; i < len(lines); i += 2 {
		if i+1 < len(lines) {
			if !isAfterPanic {
				if strings.Contains(lines[i+1], "panic") {
					isAfterPanic = true
				}
				continue
			}

			frames = append(frames, Frame{
				Location: code.NewLocation(lines[i+1]),
				Code:     code.NewCode(lines[i]),
			})
		}
	}

	return
}
