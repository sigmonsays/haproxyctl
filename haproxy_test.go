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

	info, err := ctl.Info()
	if err != nil {
		fmt.Printf("Info %s\n", err)
	}

	fmt.Printf("info %+v\n", info)
}
