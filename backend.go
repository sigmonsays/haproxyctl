package haproxyctl

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"
)

type Backend struct {
	name    string
	servers []*Server
}

func NewBackend(name string) *Backend {
	b := &Backend{
		name:    name,
		servers: make([]*Server, 0),
	}
	return b
}

func (b *Backend) Name() string {
	return b.name
}

// returns a list of backends
func (c *ControlSocket) Backends() ([]*Backend, error) {
	cmd := []byte("show backend\n")

	buf, err := c.roundTrip(cmd)
	if err != nil {
		return nil, err
	}

	if buf.Len() == 0 {
		return nil, fmt.Errorf("zero bytss received")
	}

	res := make([]*Backend, 0)

	var line []byte
	rdr := bufio.NewReader(buf)

Reader:
	for {
		line, err = rdr.ReadBytes(NL)
		if err != nil {
			break Reader
		}

		if line[0] == '#' {
			continue
		}
		if bytes.Compare(line, []byte{NL}) == 0 {
			break
		}

		name := strings.Trim(string(line), " \n")
		be := NewBackend(name)
		res = append(res, be)
	}
	if err == io.EOF {
		err = nil
	}
	if err != nil {
		return nil, err
	}

	return res, nil
}
