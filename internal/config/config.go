//	{
//	  "db_url": "connection_string_goes_here",
//	  "current_user_name": "username_goes_here"
//	}
package config

import (
	"encoding/json"
	"os"
	"path"
)

type Config struct {
	DBUrl           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func (c *Config) SetUser(user string) error {
	c.CurrentUserName = user
	err := write(*c)
	if err != nil {
		return err
	}
	return nil
}

const configFileName = ".gatorconfig.json"

func getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return path.Join(homeDir, configFileName), nil
}

func Read() (cfg *Config, err error) {
	cfg = new(Config)
	configPath, err := getConfigFilePath()
	if err != nil {
		return cfg, err
	}

	contents, err := os.ReadFile(configPath)
	if err != nil {
		return cfg, err
	}

	err = json.Unmarshal(contents, cfg)
	if err != nil {
		return cfg, err
	}

	return cfg, nil
}

func write(cfg Config) error {
	configPath, err := getConfigFilePath()
	if err != nil {
		return err
	}

	enc, err := json.Marshal(cfg)
	if err != nil {
		return err
	}

	err = os.WriteFile(configPath, enc, os.ModeCharDevice)
	if err != nil {
		return err
	}

	return nil
}
