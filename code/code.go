package code

import (
	"regexp"
	"strings"

	"github.com/gookit/color"
)

type Code struct {
	Package, Func string
}

func NewCode(s string) (c Code) {
	// github.com/moodysanalytics/maers-brt-bc-atos/internal/errors.(*Test).WithPointerReceiver
	// github.com/moodysanalytics/maers-brt-bc-atos/internal/errors.Test.WithReceiver
	r := regexp.MustCompile(`(?P<package>.+\/\w+)\.(?P<func>.+)`)
	if submatch := r.FindStringSubmatch(s); submatch != nil {
		match := make(map[string]string)
		for i, name := range r.SubexpNames() {
			if i != 0 && name != "" && len(submatch) > i {
				match[name] = submatch[i]
			}
		}

		c.Package = match["package"]
		c.Func = match["func"]
		return
	}

	// created by testing.(*T).Run in goroutine 51
	s = regexp.MustCompile(`created by (.*) in goroutine.*`).ReplaceAllString(s, `$1`)

	// testing.tRunner
	// panic({0x104de9d00?, 0x104e30f70?})
	m := strings.SplitN(s, ".", 2)
	if len(m) > 1 {
		c.Package = m[0]
		c.Func = m[1]
	} else {
		c.Func = m[0]
	}

	return
}

func (c Code) String() string {
	var s string

	if c.Package != "" {
		s = color.Yellow.Sprint(c.Package) + "."
	}

	return s + color.Green.Sprint(c.Func+"()")
}
