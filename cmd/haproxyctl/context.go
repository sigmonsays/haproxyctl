package main

import (
	"golang.org/x/crypto/ssh/terminal"

	"github.com/sigmonsays/haproxyctl"
)

type Context struct {

	// the top level root command
	Root *Command

	// the terminal
	Term *terminal.Terminal

	// the completer to use (not tab, just kepress)
	AutoCompleteCallback Completer

	// tab complete command
	TabComplete Completer

	// current command
	Current *Command

	// haproxy control
	Ha *haproxyctl.ControlSocket
}

func (c *Context) SetCurrent(cmd *Command) {
	if cmd == nil {
		cmd = c.Root
		Dbg("reset current command to name=root")
	}

	Dbg("switch current command to name=%s", cmd.name)
	c.Current = cmd
	c.AutoCompleteCallback = cmd.Completer
	c.TabComplete = cmd.TabComplete
}
