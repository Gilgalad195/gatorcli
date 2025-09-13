package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	DBUrl           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("error occurred getting home directory: %v", err)
	}
	filePath := filepath.Join(homeDir, configFileName)
	return filePath, nil
}

func write(cfg *Config) error {
	filePath, err := getConfigFilePath()
	if err != nil {
		return fmt.Errorf("error occurred: %v", err)
	}
	data, err := json.MarshalIndent(*cfg, "", " ")
	if err != nil {
		return fmt.Errorf("marshal to json failed: %v", err)
	}

	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("error writing file: %v", err)
	}
	return nil
}

func Read() (Config, error) {
	var config Config
	filePath, err := getConfigFilePath()
	if err != nil {
		return config, fmt.Errorf("error occurred: %v", err)
	}
	data, err := os.ReadFile(filePath)
	if err != nil {
		return config, fmt.Errorf("error reading file: %v", err)
	}
	if err := json.Unmarshal(data, &config); err != nil {
		return config, fmt.Errorf("error unmarshaling json: %v", err)
	}
	return config, nil
}

func (c *Config) SetUser(userName string) error {
	c.CurrentUserName = userName
	err := write(c)
	if err != nil {
		return fmt.Errorf("failed to set user: %v", err)
	}
	return nil
}

// func test() error {
// 	return nil
// }
