package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Feed struct {
	ID            uuid.UUID `json:"id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	LastFetchedAt time.Time `json:"last_fetched_at"`
	Name          string    `json:"name"`
	Url           string    `json:"url"`
	UserID        uuid.UUID `json:"user_id"`
}

func handlerAggregate(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("agg command expects 1 argumnt: duration")
	}

	timeBetweenReqs, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		return fmt.Errorf("unable to parse duration %v: %v", cmd.args[0], err)
	}
	fmt.Printf("Fetching feeds every %v\n", timeBetweenReqs)
	ticker := time.NewTicker(timeBetweenReqs)
	for ; ; <-ticker.C {
		err = scrapeFeeds(s)
		if err != nil {
			return err
		}
	}
}

// func (s *state) handlerAddFeed(w http.ResponseWriter, r *http.Request) {
// 	// if len(cmd.args) != 2 {
// 	// 	return fmt.Errorf("feed command requires 2 sub args: name and url")
// 	// }

// 	type parameters struct {
// 		Name string `json:"name"`
// 		Url  string `json:"url"`
// 	}

// 	decoder := json.NewDecoder(r.Body)
// 	params := parameters{}
// 	err := decoder.Decode(&params)
// 	if err != nil {
// 		respondWithError(w, http.StatusInternalServerError, "unable to decode params", err)
// 		return
// 	}

// 	feedParams := database.CreateFeedParams{
// 		ID:        uuid.New(),
// 		CreatedAt: time.Now().UTC(),
// 		UpdatedAt: time.Now().UTC(),
// 		Name:      params.Name,
// 		Url:       params.Url,
// 		UserID:    user.ID,
// 	}

// 	feed, err := s.db.CreateFeed(context.Background(), feedParams)
// 	if err != nil {
// 		return fmt.Errorf("unable to create RSS feed: %v", err)
// 	}

// 	feedFollowParams := database.CreateFeedFollowParams{
// 		ID:        uuid.New(),
// 		CreatedAt: time.Now().UTC(),
// 		UpdatedAt: time.Now().UTC(),
// 		UserID:    user.ID,
// 		FeedID:    feed.ID,
// 	}

// 	_, err = s.db.CreateFeedFollow(context.Background(), feedFollowParams)
// 	if err != nil {
// 		return fmt.Errorf("unable to follow feed %v for user %v: %v", feed.Name, user.Name, err)
// 	}

// 	fmt.Printf("Feed created successfully with name %v at url %v\n", feed.Name, feed.Url)
// 	return nil
// }

func handlerGetFeeds(s *state, cmd command) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("unable to get feeds")
	}

	for _, feed := range feeds {
		user, err := s.db.GetUserNameByID(context.Background(), feed.UserID)
		if err != nil {
			fmt.Printf("unable to get retrieve name for feed %v\n", feed.ID)
			continue
		}
		fmt.Println("***")
		fmt.Printf("* Feed: %v\n", feed.Name)
		fmt.Printf("* URL: %v\n", feed.Url)
		fmt.Printf("* Submitted by: %v\n", user)
		fmt.Println("***")
	}

	return nil
}
