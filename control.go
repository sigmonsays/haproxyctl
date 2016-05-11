package haproxyctl

import (
	"bufio"
	"bytes"
	"log"
	"net"
	"strings"
)

var NL = byte('\n')

type Options struct {
	Dial func(network string, address string) (net.Conn, error)

	Debug bool
}

func DefaultOptions() *Options {
	return &Options{}
}

type ControlSocket struct {
	address string
	dial    func(network string, address string) (net.Conn, error)
	logfn   func(format string, v ...interface{})
}

func NewControlSocket(address string, opts *Options) (*ControlSocket, error) {
	if opts == nil {
		opts = DefaultOptions()
	}
	if opts.Dial == nil {
		opts.Dial = net.Dial
	}

	c := &ControlSocket{
		address: address,
		dial:    opts.Dial,
	}

	if opts.Debug {
		c.logfn = log.Printf
	}
	return c, nil
}

func (c *ControlSocket) log(format string, v ...interface{}) {
	if c.logfn == nil {
		return
	}
	c.logfn(format, v...)
}

func (c *ControlSocket) dialOnce() (net.Conn, error) {
	network := "tcp"
	if strings.HasPrefix(c.address, "/") {
		network = "unix"
	}
	var err error
	conn, err := c.dial(network, c.address)

	return conn, err
}

func (c *ControlSocket) getConn() (net.Conn, error) {
	// todo: we could keep our connection open
	return c.dialOnce()
}

func (c *ControlSocket) roundTrip(request []byte) (*bytes.Buffer, error) {
	conn, err := c.getConn()
	if err != nil {
		return nil, err
	}

	defer conn.Close()

	_, err = conn.Write(request)
	if err != nil {
		return nil, err
	}

	err = nil
	res := bytes.NewBuffer(nil)
	var line []byte
	rdr := bufio.NewReader(conn)

Reader:
	for {
		line, err = rdr.ReadBytes(NL)
		if err != nil {
			break Reader
		}
		_, err = res.Write(line)
		if err != nil {
			break Reader
		}
		if bytes.Compare(line, []byte{NL}) == 0 {
			break
		}
	}
	if err != nil {
		return nil, err
	}

	return res, nil

}

func (c *ControlSocket) Close() error {
	return nil
}
