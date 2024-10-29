package main

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/imeltsner/gator-api/internal/auth"
	"github.com/imeltsner/gator-api/internal/database"
)

type Post struct {
	ID          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Title       string    `json:"title"`
	Url         string    `json:"url"`
	Description string    `json:"description,omitempty"`
	PublishedAt time.Time `json:"published_at,omitempty"`
	FeedID      uuid.UUID `json:"feed_id"`
}

func (s *state) handlerBrowse(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Limit int `json:"limit"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "unable to decode params", err)
		return
	}

	authToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "unable to parse auth header", err)
		return
	}

	userID, err := auth.ValidateJWT(authToken, s.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "unable to validate jwt", err)
		return
	}

	getPostParams := database.GetPostsForUserParams{
		UserID: userID,
		Limit:  int32(params.Limit),
	}
	posts, err := s.db.GetPostsForUser(context.Background(), getPostParams)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "unable to get post", err)
		return
	}

	type response struct {
		Posts []Post `json:"posts"`
	}
	allPosts := make([]Post, len(posts))
	for i, post := range posts {
		allPosts[i] = Post{
			ID:          post.ID,
			CreatedAt:   post.CreatedAt,
			UpdatedAt:   post.UpdatedAt,
			Title:       post.Title,
			Url:         post.Url,
			Description: post.Description.String,
			PublishedAt: post.PublishedAt.Time,
			FeedID:      post.FeedID,
		}
	}

	respondWithJSON(w, http.StatusOK, response{Posts: allPosts})
}
