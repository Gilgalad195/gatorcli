package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/Gilgalad195/gatorcli/internal/config"
	"github.com/Gilgalad195/gatorcli/internal/database"
	_ "github.com/lib/pq"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Printf("error reading file: %v\n", err)
		return
	}

	var s state
	s.cfg = &cfg
	var cmds commands
	cmds.commandMap = make(map[string]func(*state, command) error)

	db, err := sql.Open("postgres", s.cfg.DBUrl)
	if err != nil {
		fmt.Printf("error opening database: %v\n", err)
		return
	}

	dbQueries := database.New(db)
	s.db = dbQueries

	if len(os.Args) < 2 {
		fmt.Println("please enter a command")
		os.Exit(1)
	}

	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerUsers)

	cmdName := os.Args[1]
	args := os.Args[2:]

	var cmd command
	cmd.name = cmdName
	cmd.args = args

	if err := cmds.run(&s, cmd); err != nil {
		fmt.Printf("error occurred running command: %v\n", err)
		os.Exit(1)
	}

}
