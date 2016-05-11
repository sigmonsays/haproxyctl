package main

import (
	"io"
)

func Commands(ctx *Context) *Command {

	root := NewCommand("")
	root.Completer = func(line string, i int, r rune) (newline string, newpos int, ok bool) {
		return completeCommands(ctx, root, line, i, r)
	}

	quit := NewCommand("quit")
	quit.Dispatch = func(line string) error {
		return io.EOF
	}
	root.Add("quit", quit)

	initShowCommand(ctx, root)

	return root
}
