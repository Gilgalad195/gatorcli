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
	"github.com/Gilgalad195/gatorcli/internal/webconn"
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
		return err
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
	log.Printf("User %s was registered.\n", user.Name)
	return nil
}

func handlerReset(s *state, cmd command) error {
	ctx := context.Background()
	err := s.db.ResetDatabase(ctx)
	if err != nil {
		fmt.Printf("error resetting table: %v\n", err)
		os.Exit(1)
	}
	return nil
}

func handlerUsers(s *state, cmd command) error {
	fmt.Println("Retrieving users...")
	ctx := context.Background()
	users, err := s.db.GetUsers(ctx)
	if err != nil {
		fmt.Printf("error getting users: %v", err)
	}
	for _, user := range users {
		if user.Name == s.cfg.CurrentUserName {
			fmt.Printf("* %s (current)\n", user.Name)
		} else {
			fmt.Printf("* %s\n", user.Name)
		}
	}
	return nil
}

func handlerAgg(s *state, cmd command) error {
	ctx := context.Background()
	rssFeed, err := webconn.FetchFeed(ctx, "https://www.wagslane.dev/index.xml")
	if err != nil {
		return err
	}
	fmt.Printf("%+v\n", rssFeed)
	return nil
}

func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.args) < 2 {
		return fmt.Errorf("this function expects 2 arguments: name and url of the feed")
	}
	ctx := context.Background()
	user, err := s.db.GetUser(ctx, s.cfg.CurrentUserName)
	if err != nil {
		return err
	}
	feedName := cmd.args[0]
	feedURL := cmd.args[1]

	newFeed := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      feedName,
		Url:       feedURL,
		UserID:    user.ID,
	}

	feed, err := s.db.CreateFeed(ctx, newFeed)
	if err != nil {
		return fmt.Errorf("unable to create feed: %v", err)
	}

	fmt.Println("New feed was created:")
	fmt.Printf("* id: %v\n", feed.ID)
	fmt.Printf("* created_id: %v\n", feed.CreatedAt)
	fmt.Printf("* updated_at: %v\n", feed.UpdatedAt)
	fmt.Printf("* name: %v\n", feed.Name)
	fmt.Printf("* url: %v\n", feed.Url)
	fmt.Printf("* user_id: %v\n", feed.UserID)

	return nil
}

func handlerFeeds(s *state, cmd command) error {
	fmt.Println("Retrieving feeds...")
	ctx := context.Background()
	feeds, err := s.db.GetFeeds(ctx)
	if err != nil {
		fmt.Printf("error getting feeds: %v", err)
	}
	for _, feed := range feeds {
		fmt.Printf("Feed Name: %v\n", feed.Name)
		fmt.Printf("* URL: %v\n", feed.Url)
		userName, err := s.db.GetUserName(ctx, feed.UserID)
		if err != nil {
			return err
		}
		fmt.Printf("* Created by: %v\n", userName)
	}
	return nil
}
