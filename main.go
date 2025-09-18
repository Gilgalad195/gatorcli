package main

import (
	"fmt"
	"os"

	"github.com/Gilgalad195/gatorcli/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Printf("error reading file: %v\n", err)
		return
	}

	var s state
	s.configPointer = &cfg
	var cmds commands
	cmds.commandMap = make(map[string]func(*state, command) error)

	cmds.register("login", handlerLogin)
	if len(os.Args) < 2 {
		fmt.Println("please enter a command")
		os.Exit(1)
	}
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
