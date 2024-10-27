package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/imeltsner/gator-api/internal/auth"
	"github.com/imeltsner/gator-api/internal/database"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
}

func (s *state) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name      string `json:"name"`
		Password  string `json:"password"`
		ExpiresIn int    `json:"expires_in_seconds"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't decode parameters", err)
		return
	}

	dbUser, err := s.db.GetUserByName(r.Context(), params.Name)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "user not found", err)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(dbUser.HashedPassword), []byte(params.Password))
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "incorrect password", err)
		return
	}

	var expirationTime time.Duration
	if params.ExpiresIn > 0 && params.ExpiresIn > 3600 {
		expirationTime = time.Hour
	} else {
		expirationTime = time.Duration(params.ExpiresIn) * time.Second
	}

	token, err := auth.MakeJWT(dbUser.ID, s.jwtSecret, expirationTime)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "unable to make jwt", err)
		return
	}

	type response struct {
		User
		Token string `json:"token"`
	}

	respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID:        dbUser.ID,
			CreatedAt: dbUser.CreatedAt,
			UpdatedAt: dbUser.UpdatedAt,
			Name:      dbUser.Name,
		},
		Token: token,
	})
}

func (s *state) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name     string `json:"name"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "unable to decode params", err)
		return
	}
	if params.Password == "" {
		respondWithError(w, http.StatusBadRequest, "password is required", err)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.MinCost)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "unable to hash password", err)
		return
	}

	user := database.CreateUserParams{
		ID:             uuid.New(),
		CreatedAt:      time.Now().UTC(),
		UpdatedAt:      time.Now().UTC(),
		Name:           params.Name,
		HashedPassword: string(hashedPassword),
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

func (s *state) handlerGetUser(w http.ResponseWriter, r *http.Request) {
	idString := r.PathValue("id")
	id, err := uuid.Parse(idString)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "unable to parse id", err)
		return
	}

	authToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "unable to parse auth header", err)
		return
	}

	authID, err := auth.ValidateJWT(authToken, s.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "unable to validate jwt", err)
		return
	}

	if authID != id {
		respondWithError(w, http.StatusUnauthorized, "mismatched id", err)
		return
	}

	user, err := s.db.GetUserByID(r.Context(), id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "user not found", err)
		return
	}

	respondWithJSON(w, http.StatusOK, User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Name:      user.Name,
	})
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

func (s *state) handlerDeleteUser(w http.ResponseWriter, r *http.Request) {
	idString := r.PathValue("id")
	id, err := uuid.Parse(idString)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "unable to parse id", err)
		return
	}

	authToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "unable to parse auth header", err)
		return
	}

	authID, err := auth.ValidateJWT(authToken, s.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "unable to validate jwt", err)
		return
	}

	if authID != id {
		respondWithError(w, http.StatusUnauthorized, "mismatched id", err)
		return
	}

	err = s.db.DeleteUser(r.Context(), id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "unable to delete user", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *state) handlerDeleteUsers(w http.ResponseWriter, r *http.Request) {
	err := s.db.DeleteUsers(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "unable to delete users", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
