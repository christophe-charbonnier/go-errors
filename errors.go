package errors

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/gookit/color"
	"github.com/pkg/errors"
	"github.com/samber/lo"
)

func New(format string, a ...any) errorT {
	return new(errors.Errorf(format, a...))
}

func WithError(err error) errorT {
	return new(errors.WithStack(err))
}

func new(err error) errorT {
	frames := getFrames(err)[1:]
	width := lo.Max(lo.Map(frames, func(f frame, _ int) int {
		return len(f.Location.String())
	}))

	stackTrace := lo.Map(frames, func(f frame, _ int) string {
		return f.String(width + 1)
	})

	return errorT{
		Message:    err.Error(),
		StackTrace: stackTrace,
		frames:     frames,
	}
}

type errorT struct {
	Message    string
	StackTrace []string

	frames []frame
}

func (e errorT) String() string {
	return e.Error()
}

func (e errorT) Error() string {
	width := lo.Max(lo.Map(e.frames, func(f frame, _ int) int {
		return len(f.Location.String())
	}))

	message := lo.Map(strings.Split("ERROR: "+e.Message, "\n"), func(s string, _ int) string {
		return color.LightRed.Sprintf(s)
	})

	stackTrace := lo.Map(e.frames, func(f frame, _ int) string {
		return f.String(width + 1)
	})

	return fmt.Sprintf("%v\n%v",
		strings.Join(message, "\n"),
		strings.Join(stackTrace, "\n"),
	)
}

func (e errorT) Frames() []frame {
	return e.frames
}

type location struct {
	Dir, File, Line string
}

func newLocation(s string) (l location) {
	a := strings.Split(s, ":")
	file := a[0]
	l.Line = a[1]

	a = strings.SplitAfter(file, "/")
	l.Dir = strings.Join(a[0:len(a)-1], "")
	l.File = a[len(a)-1]

	return
}

func (l location) String() string {
	return color.Gray.Sprint(l.Dir) + color.LightCyan.Sprint(l.File) + color.Gray.Sprint(":") + color.Cyan.Sprint(l.Line)
}

type code struct {
	Package, Func string
}

func newCode(s string) (c code) {
	functionWithReceiver := regexp.MustCompile(`(.+)\.(\(.+)`)
	functionSimple := regexp.MustCompile(`(.+)\.(.+)`)

	a := functionWithReceiver.FindAllStringSubmatch(s, -1)
	if len(a) == 0 {
		a = functionSimple.FindAllStringSubmatch(s, -1)
	}

	if len(a) != 1 || len(a[0]) != 3 {
		return
	}

	c.Package = a[0][1]
	c.Func = a[0][2]

	return
}

func (c code) String() string {
	return color.Yellow.Sprint(c.Package) + "." + color.LightYellow.Sprint(c.Func+"()")
}

type frame struct {
	Location location
	Code     code
}

func newFrame(s string) (f frame) {
	a := regexp.MustCompile(`\s+`).Split(s, -1)
	f.Code = newCode(a[0])
	f.Location = newLocation(a[1])

	return
}

func (f frame) String(colWidth int) string {
	return fmt.Sprintf("%-*v", colWidth, f.Location.String()) + f.Code.String()
}

func getFrames(err error) (frames []frame) {
	type stackTracer interface {
		StackTrace() errors.StackTrace
	}

	if err, ok := err.(stackTracer); ok {
		for _, f := range err.StackTrace() {
			frames = append(frames, newFrame(fmt.Sprintf("%+v", f)))
		}
	}

	return
}
