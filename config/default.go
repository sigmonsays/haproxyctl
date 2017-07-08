package config

// based on https://cbonte.github.io/haproxy-dconv/1.6/configuration.html

func NewDefault() *Default {
	return &Default{
		Values: make([]Line, 0),
	}
}

type Default struct {

	// blindly store all defaults for now
	Values []Line
}

// parse the lines in a default
func (me *Parser) defaultSection(line Line, state *state) error {
	log.Tracef("parse %s", line)
	me.Config.Default.Values = append(me.Config.Default.Values, line)

	return nil
}
