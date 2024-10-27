package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/imeltsner/gator-api/internal/config"
	"github.com/imeltsner/gator-api/internal/database"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

type state struct {
	cfg       *config.Config
	db        *database.Queries
	jwtSecret string
}

func main() {
	// Check cmd line args
	//args := os.Args
	// if len(args) < 2 {
	// 	fmt.Printf("program expects at least 1 arg\n")
	// 	os.Exit(1)
	// }

	// Read config file
	// cfg, err := config.Read()
	// if err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }
	cfg := config.Config{}

	// Load environment variables
	godotenv.Load()
	dbString := os.Getenv("DB_CONNECTION")

	// Connect to db
	db, err := sql.Open("postgres", dbString)
	if err != nil {
		fmt.Printf("unable to connect to db: %v", err)
		os.Exit(1)
	}
	dbQueries := database.New(db)

	s := state{
		cfg:       &cfg,
		db:        dbQueries,
		jwtSecret: os.Getenv("JWT_SECRET"),
	}

	// Create http server
	port := os.Getenv("PORT")
	mux := http.NewServeMux()
	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	// Register routes
	mux.HandleFunc("POST /api/login", s.handlerLogin)
	mux.HandleFunc("POST /api/users", s.handlerCreateUser)
	mux.HandleFunc("GET /api/users/{id}", s.handlerGetUser)
	mux.HandleFunc("GET /api/users", s.handlerGetUsers)
	mux.HandleFunc("DELETE /api/users/{id}", s.handlerDeleteUser)
	mux.HandleFunc("DELETE /admin/reset", s.handlerDeleteUsers)

	// Register commands
	cmds := commands{
		cmds: map[string]func(*state, command) error{},
	}
	//cmds.register("login", handlerLogin)
	//cmds.register("register", handlerRegister)
	//cmds.register("reset", handlerReset)
	//cmds.register("users", handlerUsers)
	cmds.register("agg", handlerAggregate)
	cmds.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	cmds.register("feeds", handlerGetFeeds)
	cmds.register("follow", middlewareLoggedIn(handlerFollow))
	cmds.register("following", middlewareLoggedIn(handlerFollowing))
	cmds.register("unfollow", middlewareLoggedIn(handlerUnfollow))
	cmds.register("browse", middlewareLoggedIn(handlerBrowse))

	// Parse cmd line args
	// cmdName := args[1]
	// cmdSubArgs := args[2:]
	// cmd := command{
	// 	name: cmdName,
	// 	args: cmdSubArgs,
	// }

	// Run cmd
	// err = cmds.run(&s, cmd)
	// if err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }

	// Start server
	log.Printf("Serving on port: %s\n", port)
	log.Fatal(server.ListenAndServe())
}
