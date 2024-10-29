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

type FeedFollow struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserID    uuid.UUID `json:"user_id"`
	FeedID    uuid.UUID `json:"feed_id"`
}

func (s *state) handlerFollow(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Url string `json:"url"`
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

	feed, err := s.db.GetFeedByURL(context.Background(), params.Url)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "feed not found", err)
		return
	}

	feedFollowParams := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    userID,
		FeedID:    feed.ID,
	}

	feedFollow, err := s.db.CreateFeedFollow(context.Background(), feedFollowParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "unable to create feed follow entry", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, FeedFollow{
		ID:        feedFollow.ID,
		CreatedAt: feedFollow.CreatedAt,
		UpdatedAt: feedFollow.UpdatedAt,
		UserID:    userID,
		FeedID:    feed.ID,
	})
}

func (s *state) handlerFollowing(w http.ResponseWriter, r *http.Request) {
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

	feeds, err := s.db.GetFeedFollowsForUser(context.Background(), userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "unable to get follows", err)
		return
	}

	type feedFollowForUser struct {
		feedFollow FeedFollow
		Title      string `json:"title"`
		PostedBy   string `json:"posted_by"`
	}
	type response struct {
		FeedsFollowed []feedFollowForUser `json:"feeds_followed"`
	}
	feedsFollowed := make([]feedFollowForUser, len(feeds))
	for i, feed := range feeds {
		feedsFollowed[i] = feedFollowForUser{
			feedFollow: FeedFollow{
				ID:        feed.ID,
				CreatedAt: feed.CreatedAt,
				UpdatedAt: feed.UpdatedAt,
				UserID:    feed.UserID,
				FeedID:    feed.FeedID,
			},
			Title:    feed.FeedTitle,
			PostedBy: feed.UserName,
		}
	}

	respondWithJSON(w, http.StatusOK, response{FeedsFollowed: feedsFollowed})
}

func (s *state) handlerUnfollow(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Url string `json:"url"`
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

	feed, err := s.db.GetFeedByURL(context.Background(), params.Url)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "feed not found", err)
		return
	}

	deleteParams := database.DeleteFeedFollowParams{
		UserID: userID,
		FeedID: feed.ID,
	}
	err = s.db.DeleteFeedFollow(context.Background(), deleteParams)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "unable to unfollow", err)
	}

	w.WriteHeader(http.StatusNoContent)
}
