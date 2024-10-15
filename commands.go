package main

import (
	"fmt"

	"github.com/imeltsner/gator/internal/config"
)

type state struct {
	cfg *config.Config
}

type command struct {
	name string
	args []string
}

type commands struct {
	cmds map[string]func(*state, command) error
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("login command expects 1 argument")
	}

	err := s.cfg.SetUser(cmd.args[0])
	if err != nil {
		return err
	}

	fmt.Printf("User has been set to %v\n", cmd.args[0])
	return nil
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
	}

	return fmt.Errorf("command %v not found", cmd.name)
}
