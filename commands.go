package main

import (
	"fmt"

	"github.com/Gilgalad195/gatorcli/internal/config"
)

type state struct {
	configPointer *config.Config
}

type command struct {
	name string
	args []string
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("the login handler expects a single argument, the username")
	}
	err := s.configPointer.SetUser(cmd.args[0])
	if err != nil {
		return fmt.Errorf("error occurred setting user: %v", err)
	}
	fmt.Printf("terminal user has been set to %s", cmd.args[0])
	return nil
}
