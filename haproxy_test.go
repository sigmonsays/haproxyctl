package haproxyctl

import (
	"fmt"
	"testing"
)

func Test_Sanity(t *testing.T) {

	sockpath := "/var/run/haproxy.sock"

	opts := DefaultOptions()
	opts.Debug = true

	ctl, err := NewControlSocket(sockpath, opts)
	if err != nil {
		fmt.Printf("connect error %s\n", err)
		return
	}

	defer ctl.Close()

	ok, err := ctl.Ping()
	if ok == true {
		fmt.Printf("ping ok\n")
	}

	info, err := ctl.Info()
	if err == nil {
		fmt.Printf("haproxy version %s\n", info.Version)
	} else {
		fmt.Printf("info error %s\n", err)
	}
	/*
		stats, err := ctl.Stat()
		if err != nil {
			fmt.Printf("stat %s\n", err)
		}

		for _, st := range stats {
			fmt.Printf("stat %+v\n", st)
		}
	*/

	/*
		state, err := ctl.ServerStateAll()
		if err != nil {
			fmt.Printf("state %s\n", err)
		}

		for _, st := range state {
			fmt.Printf("state %+v\n", st)
		}
	*/

	backends, err := ctl.Backends()
	for _, b := range backends {
		fmt.Printf("backend %+v\n", b.Name())

		// get the servers for each backend

		servers, err := ctl.Servers(b)
		if err != nil {
			fmt.Printf("error: %s\n", err)
			continue
		}
		for _, s := range servers {
			fmt.Printf(" - %s\n", s.String())
		}
	}
}
