package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/imeltsner/gator/internal/config"
	"github.com/imeltsner/gator/internal/database"

	_ "github.com/lib/pq"
)

type state struct {
	cfg *config.Config
	db  *database.Queries
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

	db, err := sql.Open("postgres", cfg.DBURL)
	if err != nil {
		fmt.Printf("unable to connect to db: %v", err)
	}
	dbQueries := database.New(db)

	s := state{
		cfg: &cfg,
		db:  dbQueries,
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
