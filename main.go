package main

import (
	"fmt"
	"os"

	"github.com/imeltsner/gator/internal/config"
)

type state struct {
	cfg *config.Config
}

func main() {
	args := os.Args
	if len(args) < 2 {
		fmt.Printf("program expects at least 1 arg\n")
		os.Exit(1)
	}

	cfg, err := config.Read()
	if err != nil {
		os.Exit(1)
	}

	s := state{
		cfg: &cfg,
	}

	cmds := commands{
		cmds: map[string]func(*state, command) error{},
	}
	cmds.register("login", handlerLogin)

	cmdName := args[1]
	var cmdSubArgs []string
	if len(args) > 2 {
		cmdSubArgs = args[2:]
	}

	cmd := command{
		name: cmdName,
		args: cmdSubArgs,
	}

	err = cmds.run(&s, cmd)
	if err != nil {
		os.Exit(1)
	}
}
