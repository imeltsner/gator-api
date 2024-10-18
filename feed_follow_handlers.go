package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/imeltsner/gator/internal/database"
)

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("follow cmd expects 1 arg: url")
	}

	feed, err := s.db.GetFeedByURL(context.Background(), cmd.args[0])
	if err != nil {
		return fmt.Errorf("unable to get feeed by url: %v", err)
	}

	feedFollowParams := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}

	feedFollow, err := s.db.CreateFeedFollow(context.Background(), feedFollowParams)
	if err != nil {
		return fmt.Errorf("unable to create feed follow row: %v", err)
	}

	fmt.Printf("User %v successfully followed feed %v\n", feedFollow.UserName, feedFollow.FeedName)
	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {
	feeds, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("unable to get feeds user %v is following: %v", s.cfg.CurrentUsername, err)
	}

	fmt.Printf("User %v is following:\n", s.cfg.CurrentUsername)
	for _, feed := range feeds {
		fmt.Printf("* %v\n", feed.FeedName)
	}

	return nil
}
