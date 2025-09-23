package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Gilgalad195/gatorcli/internal/config"
	"github.com/Gilgalad195/gatorcli/internal/database"
	"github.com/google/uuid"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}

type command struct {
	name string
	args []string
}

type commands struct {
	commandMap map[string]func(*state, command) error
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

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("the login handler expects a single argument: the username")
	}

	username := cmd.args[0]
	ctx := context.Background()

	_, err := s.db.GetUser(ctx, username)
	if err == sql.ErrNoRows {
		fmt.Printf("user doesn't exist: %v", err)
		os.Exit(1)
	} else if err != nil {
		return fmt.Errorf("error checking user: %v", err)
	}

	if err := s.cfg.SetUser(username); err != nil {
		return fmt.Errorf("error occurred setting user: %v", err)
	}
	fmt.Printf("terminal user has been set to %s\n", cmd.args[0])
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("the register handler expects a single argument: the username")
	}
	username := cmd.args[0]

	newUser := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      username,
	}

	ctx := context.Background()

	_, err := s.db.GetUser(ctx, username)
	if err == nil {
		fmt.Println("User already exists")
		os.Exit(1)
	}
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	user, err := s.db.CreateUser(ctx, newUser)
	if err != nil {
		return fmt.Errorf("unable to create user: %v", err)
	}
	s.cfg.SetUser(user.Name)
	fmt.Printf("User %s was registered:\n", user.Name)
	log.Printf("- %v\n", user.ID)
	log.Printf("- %v\n", user.CreatedAt)
	log.Printf("- %v\n", user.UpdatedAt)
	log.Printf("- %s\n", user.Name)
	return nil
}
