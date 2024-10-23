package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/imeltsner/gator-api/internal/database"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
}

// TODO: add password
func (s *state) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name string `json:"username"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 500, "couldn't decode parameters", err)
		return
	}

	dbUser, err := s.db.GetUser(r.Context(), params.Name)
	if err != nil {
		respondWithError(w, 404, "user not found", err)
		return
	}

	respondWithJSON(w, 200, User{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Name:      dbUser.Name,
	})

	// err = s.cfg.SetUser(dbUser.Name)
	// if err != nil {
	// 	return err
	// }
}

func (s *state) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name string `json:"name"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 500, "unable to decode params", err)
		return
	}

	user := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
	}

	dbUser, err := s.db.CreateUser(r.Context(), user)
	if err != nil && strings.Contains(err.Error(), "duplicate key value") {
		respondWithError(w, 409, "name already exists in db", err)
		return
	} else if err != nil {
		respondWithError(w, 500, "unable to create user", err)
		return
	}

	log.Printf("DB user with name %v and ID %v created at %v\n", dbUser.Name, dbUser.ID, dbUser.CreatedAt)
	respondWithJSON(w, 201, User{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Name:      dbUser.Name,
	})

	// err = s.cfg.SetUser(cmd.args[0])
	// if err != nil {
	// 	return fmt.Errorf("unable to set user: %v", err)
	// }
}

func handlerReset(s *state, _ command) error {
	err := s.db.DeleteUsers(context.Background())
	if err != nil {
		return fmt.Errorf("unable to delete users: %v", err)
	}

	fmt.Println("Users successfully deleted")
	return nil
}

func handlerUsers(s *state, _ command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("unable to get users %v", err)
	}

	for _, user := range users {
		if user.Name == s.cfg.CurrentUsername {
			fmt.Printf("* %v (current)\n", user.Name)
		} else {
			fmt.Printf("* %v\n", user.Name)
		}
	}

	return nil
}
