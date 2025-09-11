package main

import (
	"fmt"

	"github.com/Gilgalad195/gatorcli/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Printf("error reading file: %v", err)
		return
	}
	cfg.SetUser("Stephen")

	cfgUpdated, err := config.Read()
	if err != nil {
		fmt.Printf("error reading file: %v", err)
		return
	}
	fmt.Printf("db_url: %s\n", cfgUpdated.DBUrl)
	fmt.Printf("current_user_name: %s\n", cfgUpdated.CurrentUserName)
}
