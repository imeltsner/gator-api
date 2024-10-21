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
	// Check cmd line args
	args := os.Args
	if len(args) < 2 {
		fmt.Printf("program expects at least 1 arg\n")
		os.Exit(1)
	}

	// Read config file
	cfg, err := config.Read()
	if err != nil {
		os.Exit(1)
	}

	// Connect to db
	db, err := sql.Open("postgres", cfg.DBURL)
	if err != nil {
		fmt.Printf("unable to connect to db: %v", err)
	}
	dbQueries := database.New(db)

	s := state{
		cfg: &cfg,
		db:  dbQueries,
	}

	// Register commands
	cmds := commands{
		cmds: map[string]func(*state, command) error{},
	}
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerUsers)
	cmds.register("agg", handlerAggregate)
	cmds.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	cmds.register("feeds", handlerGetFeeds)
	cmds.register("follow", middlewareLoggedIn(handlerFollow))
	cmds.register("following", middlewareLoggedIn(handlerFollowing))
	cmds.register("unfollow", middlewareLoggedIn(handlerUnfollow))
	cmds.register("browse", middlewareLoggedIn(handlerBrowse))

	// Parse cmd line args
	cmdName := args[1]
	cmdSubArgs := args[2:]
	cmd := command{
		name: cmdName,
		args: cmdSubArgs,
	}

	// Run cmd
	err = cmds.run(&s, cmd)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
