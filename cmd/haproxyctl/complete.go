package main

import (
	"fmt"
	"strings"
)

func completeCommands(ctx *Context, root *Command, line string, i int, r rune) (newline string, newpos int, ok bool) {
	if line == "" {
		ctx.SetCurrent(nil)
	}

	sline := line

	// only add printables into the command we're splitting
	if r >= 32 {
		sline += string(r)
	}
	tmp := strings.Split(sline, " ")
	l := len(tmp)
	prefix := tmp[l-1]

	// do not complete the command if we have no prefix
	if prefix == "" {
		Dbg("empty prefix")
		return "", 0, false
	}

	found := make([]*Command, 0)
	for cmd, c := range root.commands {
		Dbg("test match prefix=%s cmd=%s", prefix, cmd)
		if strings.HasPrefix(cmd, prefix) {
			found = append(found, c)
		}
	}

	Dbg("completeCommands %q line=%q i=%d r=%q/%d (#found=%d)\n", root.Name(), line, i, string(r), r, len(found))
	if len(found) == 0 {
		return "", 0, false

	} else if len(found) == 1 {
		ctx.SetCurrent(found[0])
	} else {

		if r == HT { // tab
			names := make([]string, 0)
			for _, f := range found {
				names = append(names, f.Name())
			}
			fmt.Printf("arbitrary command: try %s\n", strings.Join(names, ", "))

		}
	}

	return "", 0, false

}
