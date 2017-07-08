package config

import "fmt"

// based on https://cbonte.github.io/haproxy-dconv/1.6/configuration.html

func NewGlobal() *Global {
	return &Global{
		Values: make(map[string]Line, 0),
	}
}

type Global struct {

	// blindly parse all globals for now
	Values map[string]Line
}

// parse the lines in a global
func (me *Parser) globalSection(line Line, state *state) error {
	log.Tracef("line %#v", line)
	name := line.Command()
	me.Config.Global.Values[name] = line

	return nil
}
func (me *Global) GetParam(name string) (Line, error) {
	line, found := me.Values[name]
	if found == false {
		return nil, fmt.Errorf("%s not found", name)
	}
	return line, nil
}
