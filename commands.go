package main

import (
	"fmt"
)

type command struct {
	name string
	args []string
}

type commands struct {
	cmds map[string]func(*state, command) error
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.cmds[name] = f
}

func (c *commands) run(s *state, cmd command) error {
	if handler, ok := c.cmds[cmd.name]; ok {
		err := handler(s, cmd)
		if err != nil {
			return err
		}
		return nil
	}

	return fmt.Errorf("command %v not found", cmd.name)
}
