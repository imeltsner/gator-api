package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/imeltsner/gator/internal/database"
)

func handlerAggregate(s *state, cmd command) error {
	feed, err := fetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return err
	}

	fmt.Println(feed)
	return nil
}

func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.args) != 2 {
		return fmt.Errorf("feed command requires 2 sub args: name and url")
	}

	currentUser, err := s.db.GetUser(context.Background(), s.cfg.CurrentUsername)
	if err != nil {
		return fmt.Errorf("unable to get user from db: %v", err)
	}

	feedParams := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      cmd.args[0],
		Url:       cmd.args[1],
		UserID:    currentUser.ID,
	}

	feed, err := s.db.CreateFeed(context.Background(), feedParams)
	if err != nil {
		return fmt.Errorf("unable to create RSS feed: %v", err)
	}

	fmt.Printf("Feed created successfully with name %v at url %v\n", feed.Name, feed.Url)
	return nil
}

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
