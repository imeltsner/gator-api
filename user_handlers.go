package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
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

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("register command expects 1 argument")
	}

	user := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      cmd.args[0],
	}

	dbUser, err := s.db.CreateUser(context.Background(), user)
	if err != nil {
		return fmt.Errorf("unable to create db user: %v", err)
	}

	err = s.cfg.SetUser(cmd.args[0])
	if err != nil {
		return fmt.Errorf("unable to set user: %v", err)
	}

	fmt.Printf("DB user with name %v and ID %v created at %v\n", dbUser.Name, dbUser.ID, dbUser.CreatedAt)
	return nil
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
