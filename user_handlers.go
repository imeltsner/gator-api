package main

import (
	"encoding/json"
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
		Name string `json:"name"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't decode parameters", err)
		return
	}

	dbUser, err := s.db.GetUser(r.Context(), params.Name)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "user not found", err)
		return
	}

	respondWithJSON(w, http.StatusOK, User{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Name:      dbUser.Name,
	})
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
		respondWithError(w, http.StatusConflict, "name already exists in db", err)
		return
	} else if err != nil {
		respondWithError(w, http.StatusInternalServerError, "unable to create user", err)
		return
	}

	log.Printf("DB user with name %v and ID %v created at %v\n", dbUser.Name, dbUser.ID, dbUser.CreatedAt)
	respondWithJSON(w, http.StatusCreated, User{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Name:      dbUser.Name,
	})
}

func (s *state) handlerDeleteUsers(w http.ResponseWriter, r *http.Request) {
	err := s.db.DeleteUsers(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "unable to delete users", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *state) handlerGetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := s.db.GetUsers(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "unable to get users", err)
		return
	}

	type response struct {
		Users []User `json:"users"`
	}

	allUsers := make([]User, len(users))
	for i, user := range users {
		allUsers[i] = User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Name:      user.Name,
		}
	}

	respondWithJSON(w, http.StatusOK, response{
		Users: allUsers,
	})
}
