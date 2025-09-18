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

type commands struct {
	commandMap map[string]func(*state, command) error
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("the login handler expects a single argument: the username")
	}
	err := s.configPointer.SetUser(cmd.args[0])
	if err != nil {
		return fmt.Errorf("error occurred setting user: %v", err)
	}
	fmt.Printf("terminal user has been set to %s\n", cmd.args[0])
	return nil
}

func (c *commands) run(s *state, cmd command) error {
	f := c.commandMap[cmd.name]
	if f == nil {
		return fmt.Errorf("command not recognized")
	}
	err := f(s, cmd)
	if err != nil {
		return fmt.Errorf("error occurred running command: %v", err)
	}
	return nil
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.commandMap[name] = f
}
