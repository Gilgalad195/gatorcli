package main

import (
	"github.com/Gilgalad195/gatorcli/internal/database"
)

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return nil
}
