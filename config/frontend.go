package config

import "strings"

// based on https://cbonte.github.io/haproxy-dconv/1.6/configuration.html

func NewFrontend() *Frontend {
	b := &Frontend{
		Values: make(map[string]Line, 0),
		Bind:   NewBind(),
	}
	return b
}

type Frontend struct {
	Name string

	DefaultBackend string

	Bind *Bind

	// blindly parse all defaults for now
	Values map[string]Line
}

func NewBind() *Bind {
	b := &Bind{}
	return b
}

type Bind struct {
	Address   string
	PortRange string
	Path      string
	Params    Line
}

// parse the lines in a default
func (me *Parser) frontendSection(line Line, state *state) error {
	log.Tracef("line %#v", line)

	cmd := line.Command()

	switch cmd {
	case "bind":

		bind := state.currentFrontend.Bind
		address := line.Arg(1)
		if strings.HasPrefix(address, "/") {
			bind.Path = address // its a unix address
		} else {

			// its a <address>:<port_range> address
			i := strings.LastIndex(address, ":")
			if i > 0 {
				bind.Address = address[i:]
				// todo: bind.PortRange
			}

		}
		bind.Params = line[2:]

	case "default_backend":
		state.currentFrontend.DefaultBackend = line.Arg(1)
	default:
		log.Tracef("unsupported %s", line)
	}

	return nil
}
