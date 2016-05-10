package haproxy

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
		fmt.Printf("Yippie\n")
	}

	stats, err := ctl.Stat()
	if err != nil {
		fmt.Printf("stat %s\n", err)
	}

	for _, st := range stats {
		fmt.Printf("stat %+v\n", st)
	}
}
