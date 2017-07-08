package config

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

func Parse(filename string) (*Parser, error) {
	p := NewParser(filename)
	err := p.Parse()
	return p, err
}

func NewParser(filename string) *Parser {
	p := &Parser{
		filename: filename,
		Config:   NewConfig(),
	}
	return p
}

type Parser struct {
	filename string

	Config *Config
}

type state struct {
	currentSection string
	parsefunc      func(line Line, state *state) error

	currentBackend  *Backend
	currentServer   *Server
	currentFrontend *Frontend
}

func (me *Parser) Parse() error {
	log.Tracef("parsing %s", me.filename)

	f, err := os.Open(me.filename)
	if err != nil {
		return err
	}
	defer f.Close()
	bf := bufio.NewReader(f)

	state := &state{}
	for {
		line, err := bf.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Warnf("read bytes: %s", err)
			break
		}

		sline := strings.Trim(string(line), "\n \t")
		if len(sline) == 0 {
			continue
		}
		if strings.HasPrefix(sline, "#") {
			continue
		}

		pline, err := me.parseLine(sline)
		if err != nil {
			log.Warnf("parseLine %q: %s", sline, err)
			break
		}

		// log.Tracef("%#v", pline)

		if pline.Command() == "global" {
			state.currentSection = "global"
			state.parsefunc = me.globalSection
			continue

		} else if pline.Command() == "defaults" {
			state.currentSection = "defaults"
			state.parsefunc = me.defaultSection
			continue

		} else if pline.Command() == "frontend" {
			state.currentSection = "frontend"
			state.parsefunc = me.frontendSection
			state.currentFrontend = NewFrontend()
			state.currentFrontend.Name = pline.Arg(1)
			log.Tracef("new frontend %+v", state.currentFrontend)
			me.Config.Frontend = append(me.Config.Frontend, state.currentFrontend)
			continue

		} else if pline.Command() == "backend" {
			state.currentSection = "backend"
			state.parsefunc = me.backendSection
			state.currentBackend = NewBackend()
			state.currentBackend.Name = pline.Arg(1)
			log.Tracef("new backend %+v", state.currentBackend)
			me.Config.Backend = append(me.Config.Backend, state.currentBackend)
			continue
		}

		if state.parsefunc != nil {
			err = state.parsefunc(pline, state)
			if err != nil {
				log.Warnf("parseLine %q: %s", sline, err)
				break
			}
		}
	}

	return nil
}

var (
	space       = ' '
	tab         = '\t'
	singleQuote = '\''
	doubleQuote = '"'
)

func (me *Parser) parseLine(line string) (Line, error) {
	ret := make([]string, 0)

	token := ""

	for i := 0; i < len(line); i++ {

		r := rune(line[i])
		switch r {
		case space, tab:
			if len(token) > 0 {
				ret = append(ret, token)
			}
			token = ""
		case singleQuote:
			// keep consuming until next quote (ending)
			endFound := false
			for j := i + 1; j < len(line); j++ {
				if rune(line[j]) == singleQuote {
					endFound = true
					break
				}
			}
			if endFound == false {
				return ret, fmt.Errorf("no end quote found")
			}

		default:
			token += string(line[i])
		}
	}

	if len(token) > 0 {
		ret = append(ret, token)
	}

	return Line(ret), nil
}

// represents a tokenized configuration line
type Line []string

// returns the command name
func (me Line) Command() string {
	return me[0]
}
func (me Line) Len() int {
	return len(me)
}
func (me Line) Arg(i int) string {
	l := len(me)
	if i > l-1 {
		return ""
	}
	return me[i]
}
