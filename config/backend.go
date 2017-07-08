package config

// based on https://cbonte.github.io/haproxy-dconv/1.6/configuration.html

func NewBackend() *Backend {
	b := &Backend{
		Values:  make(map[string]Line, 0),
		Servers: make([]*Server, 0),
	}
	return b
}

type Backend struct {

	// name of the backend
	Name string

	// blindly parse all defaults for now
	Values map[string]Line

	Servers []*Server
}

func NewServer() *Server {
	s := &Server{}
	return s
}

type Server struct {
	Name    string
	Address string

	Params Line
}

// parse the lines in a default
func (me *Parser) backendSection(line Line, state *state) error {
	log.Tracef("parse %s", line)

	cmd := line.Command()

	if cmd == "server" {
		state.currentServer = NewServer()
		state.currentBackend.Servers = append(state.currentBackend.Servers, state.currentServer)

		state.currentServer.Name = line.Arg(1)
		state.currentServer.Address = line.Arg(2)
		state.currentServer.Params = line[3:]

		log.Tracef("new server %+v", state.currentServer)

	}

	return nil
}
