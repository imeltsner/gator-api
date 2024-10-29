package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/imeltsner/gator-api/internal/database"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

type state struct {
	db        *database.Queries
	jwtSecret string
}

func main() {
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

	// Register user routes
	mux.HandleFunc("POST /api/login", s.handlerLogin)
	mux.HandleFunc("POST /api/users", s.handlerCreateUser)
	mux.HandleFunc("GET /api/users/{id}", s.handlerGetUser) // authenticated
	mux.HandleFunc("GET /api/users", s.handlerGetUsers)
	mux.HandleFunc("DELETE /api/users/{id}", s.handlerDeleteUser) // authenticated
	mux.HandleFunc("DELETE /admin/reset", s.handlerDeleteUsers)

	// Register feed routes
	mux.HandleFunc("POST /api/feeds", s.handlerAddFeed) // authenticated
	mux.HandleFunc("GET /api/feeds/{id}", s.handlerGetFeed)
	mux.HandleFunc("GET /api/feeds", s.handlerGetFeeds)
	mux.HandleFunc("POST /api/agg", s.handlerAggregate)

	// Register follow routes
	mux.HandleFunc("POST /api/follows", s.handlerFollow)     // authenticated
	mux.HandleFunc("GET /api/follows", s.handlerFollowing)   // authenticated
	mux.HandleFunc("DELETE /api/follows", s.handlerUnfollow) // authenticated

	// Register post routes
	mux.HandleFunc("GET /api/posts", s.handlerBrowse) // authenticated

	// Start server
	log.Printf("Serving on port: %s\n", port)
	log.Fatal(server.ListenAndServe())
}
