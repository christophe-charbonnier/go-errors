package errors

import (
	"fmt"
	"os"
	"strings"

	"github.com/gookit/color"
	"github.com/pkg/errors"
	"github.com/samber/lo"
)

func init() {
	color.Enable = os.Getenv("TERM") != ""
}

func New(format string, a ...any) Error {
	return new(errors.Errorf(format, a...))
}

func WithError(err error) Error {
	if err, ok := err.(Error); ok {
		return err
	}
	return new(errors.WithStack(err))
}

func WithRecover(recover any) Error {
	frames := FramesFromDebug()

	width := lo.Max(lo.Map(frames, func(f Frame, _ int) int {
		return len(f.Location.String())
	}))

	stackTrace := lo.Map(frames, func(f Frame, _ int) string {
		return f.String(width + 1)
	})

	return Error{
		Message:    fmt.Sprint(recover),
		StackTrace: stackTrace,
		Frames:     frames,
	}
}

func WithDebug(message ...string) Error {
	frames := FramesFromDebug()
	width := lo.Max(lo.Map(frames, func(f Frame, _ int) int {
		return len(f.Location.String())
	}))

	stackTrace := lo.Map(frames, func(f Frame, _ int) string {
		return f.String(width + 1)
	})

	return Error{
		Message:    strings.Join(message, ""),
		StackTrace: stackTrace,
		Frames:     frames,
	}
}

func new(err error) Error {
	frames := FramesFromError(err)
	if len(frames) > 0 {
		frames = frames[1:]
	}

	width := lo.Max(lo.Map(frames, func(f Frame, _ int) int {
		return len(f.Location.String())
	}))

	stackTrace := lo.Map(frames, func(f Frame, _ int) string {
		return f.String(width + 1)
	})

	return Error{
		Message:    err.Error(),
		StackTrace: stackTrace,
		Frames:     frames,
	}
}

type Error struct {
	Message    string
	StackTrace []string
	Frames     []Frame
}

func (e Error) Error() string {
	width := lo.Max(lo.Map(e.Frames, func(f Frame, _ int) int {
		return len(f.Location.String())
	}))

	message := lo.Map(strings.Split(e.Message, "\n"), func(s string, _ int) string {
		return color.LightRed.Sprintf(s)
	})

	stackTrace := lo.Map(e.Frames, func(f Frame, _ int) string {
		return f.String(width + 1)
	})

	return fmt.Sprintf("%v\n%v",
		strings.Join(message, "\n"),
		strings.Join(stackTrace, "\n"),
	)
}

func Message(err error) string {
	if err, ok := err.(Error); ok {
		return err.Message
	}
	return err.Error()
}
