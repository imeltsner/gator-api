package main

import (
	"context"
	"fmt"

	"github.com/imeltsner/gator-api/internal/database"
)

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		user, err := s.db.GetUserByName(context.Background(), s.cfg.CurrentUsername)
		if err != nil {
			return fmt.Errorf("unable to get user from db: %v", err)
		}

		return handler(s, cmd, user)
	}
}
