package config

import (
	gologging "github.com/sigmonsays/go-logging"
)

var log gologging.Logger

func init() {
	log = gologging.Register("haproxy.config", func(newlog gologging.Logger) { log = newlog })
}
