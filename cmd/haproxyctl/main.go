package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/sigmonsays/haproxyctl"
	"golang.org/x/crypto/ssh/terminal"
)

var commands = []string{
	"add acl *acl *acl_pattern",
	"add map *map *map_key *map_value",
	"clear counters",
	"clear counters all",
	"clear acl *acl",
	"clear acl *acl",
	"clear map *map",
	"help",
	// "prompt",
	"quit",
	"clear table *table",
	"del acl *acl *acl_key",
	"del map *map *map_key",
	"disable agent *backend *server",
	"disable frontend *frontend",
	"disable health *backend *server",
	"disable server *backend *server",
	"enable agent *backend *server",
	"enable frontend *frontend",
	"enable health *backend *server",
	"enable server *backend *server",
	"get map *map *map_value",
	"get acl *acl *acl_value",
	"get weight *backend *server",
	"set map *map *map_key *map_value",
	"set maxconn frontend *frontend value",
	"set maxconn server *backend *server value",
	"set maxconn global maxconn",
	"set rate-limit connections global value",
	"set rate-limit http-compression global value",
	"set rate-limit sessions global value",
	"set rate-limit ssl-sessions global value",
	"set server *backend *server addr *ip_addr",
	// agent_state = [ up | down ]
	"set server *backend *server agent *agent_state",
	// health_state = [ up | stopping | down ]
	"set server *backend *server health *health_state",
	// server_state = [ ready | drain | maint ]
	"set server *backend *server state *server_state",
	"set server *backend *server weight *server_weight",
	"set ssl ocsp-response",
	"set ssl tls-key id tlskey",
	"set table *table key *table_key data",
	"set timeout cli delay",
	"set weight *backend *server *server_weight",
	"set weight *backend *server *server_weight",
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
	"shutdown frontend *Frontend",
	"shutdown session *sid",
	"shutdown sessions *backend *server",
}

func main() {
	oldState, err := terminal.MakeRaw(0)
	if err != nil {
		panic(err)
	}
	defer terminal.Restore(0, oldState)

	sockpath := "/var/run/haproxy.sock"

	opts := haproxyctl.DefaultOptions()
	opts.Debug = true

	SetDebug("/tmp/term.log")
	Dbg("init")

	ha, err := haproxyctl.NewControlSocket(sockpath, opts)
	if err != nil {
		fmt.Printf("connect error %s\n", err)
		return
	}

	defer ha.Close()

	c := &cli{
		ha: ha,
	}

	err = c.interactive()
	if err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}
}

type cli struct {
	t    *terminal.Terminal
	ha   *haproxyctl.ControlSocket
	ctx  *Context
	root *Command
}

const (
	ETX = rune(3) // ^C
	HT  = rune(9) // <tab>
)

func (c *cli) interactive() error {
	// get haproxy info
	info, err := c.ha.Info()
	if err != nil {
		return fmt.Errorf("error getting haproxy info: %s\n", err)
	}

	// initialize the terminal
	t := terminal.NewTerminal(os.Stdin, "> ")
	c.t = t

	ctx := &Context{
		Term: t,
		Ha:   c.ha,
	}

	root := Commands(ctx)
	ctx.AutoCompleteCallback = root.Completer
	ctx.Root = root

	t.AutoCompleteCallback = func(line string, i int, r rune) (newline string, newpos int, ok bool) {
		if r == ETX { // control C
			fmt.Println("")
			return "", 0, true

		} else if r == HT { // tab
			if ctx.TabComplete == nil {
				Dbg("no tab complete for current context")
				return "", 0, false
			}

			Dbg("TabComplete line=%s i=%d r=%d", line, i, r)
			return ctx.TabComplete(line, i, r)

		} else {

			// process each keypress
			if ctx.AutoCompleteCallback == nil {
				return "", 0, false
			}

			return ctx.AutoCompleteCallback(line, i, r)
		}
	}

	fmt.Printf("connected to haproxy %s\n", info.Version)

	c.ctx = ctx
	c.root = root

	var line string
	for {

		ctx.Current = root
		ctx.AutoCompleteCallback = root.Completer

		line, err = t.ReadLine()
		if err != nil {
			Dbg("readline: %s", err)
			break
		}

		if ctx.Current == nil {
			err = c.dispatchLine(line)
			if err == io.EOF { // means quit
				fmt.Printf("bye.\n")
				break
			}
			continue
		}
		if ctx.Current.Dispatch == nil {
			Dbg("cmd %s dispatch is nil", ctx.Current.name)
			continue
		}

		err = ctx.Current.Dispatch(line)
		if err == io.EOF { // means quit
			fmt.Printf("bye.\n")
			break
		}
		if err != nil {
			fmt.Printf("%s: %s\n", ctx.Current.name, err)
		}

	}

	return nil
}

// if not a current command is set during the completer then we try to dispatch the line
func (c *cli) dispatchLine(line string) error {

	Dbg("dispatch line %q", line)

	var cmd *Command
	cmd = c.root

	tmp := strings.Split(line, " ")

	for _, name := range tmp {

		found := make([]*Command, 0)
		for _, c := range cmd.commands {
			if strings.HasPrefix(name, c.name) {
				found = append(found, c)
			}
		}

		if len(found) == 0 {
			break
		} else if len(found) == 1 {
			cmd = found[0]
			continue
		} else {
			return fmt.Errorf("ambiguous command: %s", line)
		}
	}
	Dbg("cmd name=%s", cmd.name)

	if cmd == nil || cmd == c.root {
		return fmt.Errorf("no such command: %s", line)
	}

	if cmd.Dispatch == nil {
		Dbg("cmd %s dispatch is nil", cmd.name)
		return nil
	}

	err := cmd.Dispatch(line)

	return err
}
