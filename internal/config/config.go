package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const configFileName = ".gator-apiconfig.json"

type Config struct {
	DBURL           string `json:"db_url"`
	CurrentUsername string `json:"current_user_name"`
}

func Read() (Config, error) {
	configFilePath, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}

	configFile, err := os.Open(configFilePath)
	if err != nil {
		return Config{}, fmt.Errorf("unable to open config file at path %v: %v", configFileName, err)
	}
	defer configFile.Close()

	var cfg Config
	decoder := json.NewDecoder(configFile)
	err = decoder.Decode(&cfg)
	if err != nil {
		return Config{}, fmt.Errorf("unable to decode config file: %v", err)
	}

	return cfg, nil
}

func getConfigFilePath() (string, error) {
	homePath, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("unable to get path to home dir: %v", err)
	}

	return filepath.Join(homePath, configFileName), nil
}

func (cfg *Config) SetUser(user string) error {
	cfg.CurrentUsername = user
	configFilePath, err := getConfigFilePath()
	if err != nil {
		return err
	}

	configFile, err := os.Create(configFilePath)
	if err != nil {
		return fmt.Errorf("unable to create config file: %v", err)
	}
	defer configFile.Close()

	encoder := json.NewEncoder(configFile)
	err = encoder.Encode(cfg)
	if err != nil {
		return fmt.Errorf("unable to write json to config file %v", err)
	}

	return nil
}
