package main

import (
	"context"

	"github.com/Gilgalad195/gatorcli/internal/database"
)

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		ctx := context.Background()
		username := s.cfg.CurrentUserName

		user, err := s.db.GetUser(ctx, username)
		if err != nil {
			return err
		}

		return handler(s, cmd, user)
	}
}
