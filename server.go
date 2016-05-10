package haproxy

import (
	"fmt"
)

// represents a single server in a backend
type Server struct {

	// the backend this is part of
	backend *Backend

	// record from 'show servers state'
	state *ServerState
}

func NewServer() *Server {
	s := &Server{}
	return s
}

func (s *Server) Name() string {
	return fmt.Sprintf("<Server backend=%s name=%s>", s.backend.Name(), s.state.Srv_name)
}

// return servers for a given backend
func (c *ControlSocket) Servers(backend *Backend) ([]*Server, error) {
	server_states, err := c.ServerState(backend.name)
	if err != nil {
		return nil, err
	}

	res := make([]*Server, 0)
	for _, state := range server_states {
		s := NewServer()
		s.state = state
		s.backend = backend
		res = append(res, s)
	}

	return res, nil
}
