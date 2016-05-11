package main

import (
	"fmt"

	"github.com/sigmonsays/haproxyctl"
)

/*
	"show env",
	"show errors *iid",
	"show backend",
	"show info",
	"show info typed",
	"show map *map",
	"show acl *acl",
	"show pools",
	"show servers state *backend",
	"show sess",
	"show sess *sess_id",
	"show stat *iid *type *sid typed",
	"show stat resolvers *resolver",
	"show table",
	"show table *table",
	"show tls-keys",
*/

func initShowCommand(ctx *Context, root *Command) {
	show := NewCommand("show")
	show.Dispatch = func(line string) error {
		fmt.Printf("executing show.. show what..\n")
		return nil
	}
	show.Completer = func(line string, i int, r rune) (newline string, newpos int, ok bool) {

		Dbg("show line=%v i=%d r=%v\n", line, i, r)

		return completeCommands(ctx, show, line, i, r)

	}

	show.TabComplete = func(line string, i int, r rune) (newline string, newpos int, ok bool) {
		return completeCommands(ctx, show, line, i, r)
	}

	version := show.NewCommand("version")
	version.Dispatch = func(line string) error {
		fmt.Printf("version 1.0\n")
		fmt.Printf("dbg line: %s\n", line)
		return nil
	}

	frontend := show.NewCommand("frontend")
	frontend.Dispatch = func(line string) error {
		fmt.Printf("frontend -- %q\n", line)

		frontends, _ := ctx.Ha.ShowStat(-1, haproxyctl.ObjectFrontend, -1)
		for _, f := range frontends {

			fmt.Printf("\n")
			fmt.Printf("frontend %s\n", f.Pxname)
			fmt.Printf(" %+v\n", f)
		}

		return nil
	}

	bar := show.NewCommand("bar")
	bar.Dispatch = func(line string) error {
		fmt.Printf("bar -- %q\n", line)
		return nil
	}

	root.Add("show", show)
}
