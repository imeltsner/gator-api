package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/imeltsner/gator-api/internal/database"
)

func handlerBrowse(s *state, cmd command, user database.User) error {
	var limit int
	if len(cmd.args) != 1 {
		limit = 2
	}

	limit, err := strconv.Atoi(cmd.args[0])
	if err != nil {
		return fmt.Errorf("unable to parse arg %v into int: %v", cmd.args[0], err)
	}

	getPostParams := database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  int32(limit),
	}
	posts, err := s.db.GetPostsForUser(context.Background(), getPostParams)
	if err != nil {
		return fmt.Errorf("unable to get posts for user %v: %v", user.ID, err)
	}

	fmt.Printf("Here are your %d most recent posts\n", limit)
	for _, post := range posts {
		fmt.Println("--------------------------------")
		fmt.Printf("Title: %v\n", post.Title)
		fmt.Printf("Link: %v\n", post.Url)
	}
	fmt.Println("--------------------------------")

	return nil
}
