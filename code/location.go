package code

import (
	"regexp"
	"strconv"

	"github.com/gookit/color"
)

type Location struct {
	Dir, File string
	Line      int
}

func NewLocation(s string) (l Location) {
	regexp := regexp.MustCompile(`\s*((?P<dir>.+)/)?(?P<file>.+):(?P<line>\d+)`)
	submatch := regexp.FindStringSubmatch(s)
	match := make(map[string]string)
	for i, name := range regexp.SubexpNames() {
		if i != 0 && name != "" && len(submatch) > i {
			match[name] = submatch[i]
		} else {
			match[name] = s
		}
	}

	l.Dir = match["dir"]
	l.File = match["file"]
	l.Line, _ = strconv.Atoi(match["line"])

	return
}

func (l Location) String() string {
	var s string

	if l.Dir != "" {
		s += color.Gray.Sprint(l.Dir + "/")
	}

	return s + color.LightCyan.Sprint(l.File) + color.Gray.Sprint(":") + color.Cyan.Sprint(l.Line)
}
