package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/imeltsner/gator-api/internal/auth"
	"github.com/imeltsner/gator-api/internal/database"
)

type Feed struct {
	ID            uuid.UUID `json:"id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	LastFetchedAt time.Time `json:"last_fetched_at"`
	Title         string    `json:"title"`
	Url           string    `json:"url"`
	UserID        uuid.UUID `json:"user_id"`
}

func (s *state) handlerAggregate(w http.ResponseWriter, r *http.Request) {
	err := s.scrapeFeeds()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "unable to scrape feeds", err)
	}

	w.WriteHeader(http.StatusOK)
}

func (s *state) handlerAddFeed(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Title string `json:"title"`
		Url   string `json:"url"`
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

	id, err := auth.ValidateJWT(authToken, s.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "unable to validate jwt", err)
		return
	}

	feedParams := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Title:     params.Title,
		Url:       params.Url,
		UserID:    id,
	}

	feed, err := s.db.CreateFeed(r.Context(), feedParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "unable to create feed", err)
		return
	}

	feedFollowParams := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    id,
		FeedID:    feed.ID,
	}

	_, err = s.db.CreateFeedFollow(r.Context(), feedFollowParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "unable to create feed follow entry", err)
		return
	}

	log.Printf("Feed created successfully with name %v at url %v\n", feed.Title, feed.Url)
	respondWithJSON(w, http.StatusCreated, Feed{
		ID:            feed.ID,
		CreatedAt:     feed.CreatedAt,
		UpdatedAt:     feed.UpdatedAt,
		LastFetchedAt: feed.LastFetchedAt.Time,
		Title:         feed.Title,
		Url:           feed.Url,
		UserID:        feed.UserID,
	})
}

func (s *state) handlerGetFeed(w http.ResponseWriter, r *http.Request) {
	idString := r.PathValue("id")
	feedID, err := uuid.Parse(idString)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "unable to parse feed id", err)
		return
	}

	feed, err := s.db.GetFeedByID(r.Context(), feedID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "feed not found", err)
		return
	}

	respondWithJSON(w, http.StatusOK, Feed{
		ID:            feed.ID,
		CreatedAt:     feed.CreatedAt,
		UpdatedAt:     feed.UpdatedAt,
		LastFetchedAt: feed.LastFetchedAt.Time,
		Title:         feed.Title,
		Url:           feed.Url,
		UserID:        feed.UserID,
	})
}

func (s *state) handlerGetFeeds(w http.ResponseWriter, r *http.Request) {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "unable to get feeds", err)
		return
	}

	type response struct {
		Feeds []Feed   `json:"feeds"`
		Users []string `json:"users"`
	}
	allFeeds := make([]Feed, len(feeds))
	userNames := make([]string, len(feeds))

	for i, feed := range feeds {
		user, err := s.db.GetUserNameByID(context.Background(), feed.UserID)
		if err != nil {
			respondWithError(w, http.StatusNotFound, "user not found", err)
			return
		}
		allFeeds[i] = Feed{
			ID:            feed.ID,
			CreatedAt:     feed.CreatedAt,
			UpdatedAt:     feed.UpdatedAt,
			LastFetchedAt: feed.LastFetchedAt.Time,
			Title:         feed.Title,
			Url:           feed.Url,
			UserID:        feed.UserID,
		}
		userNames[i] = user
	}

	respondWithJSON(w, http.StatusOK, response{
		Feeds: allFeeds,
		Users: userNames,
	})
}
