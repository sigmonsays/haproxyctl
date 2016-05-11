package main

func NewCommand(name string) *Command {
	return &Command{
		name:     name,
		commands: make(map[string]*Command),
	}
}

type Completer func(line string, i int, r rune) (newline string, newpos int, ok bool)

type Command struct {
	name        string
	commands    map[string]*Command
	Completer   Completer
	TabComplete Completer
	Dispatch    func(line string) error
}

func (c *Command) Name() string {
	if c.name == "" {
		return "(root)"
	}
	return c.name
}

func (c *Command) Add(name string, cmd *Command) *Command {
	c.commands[name] = cmd
	return c
}

func (c *Command) NewCommand(name string) *Command {
	newcmd := NewCommand(name)
	c.Add(name, newcmd)
	return newcmd

}
