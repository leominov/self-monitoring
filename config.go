package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

// TelegramConfig structure
type TelegramConfig struct {
	Enable    bool   `json:"enable,omitempty"`
	Token     string `json:"token,omitempty"`
	ContactID int    `json:"contactID,omitempty"`
	Debug     bool   `json:"debug,omitempty"`
}

// ConfigFile structure
type ConfigFile struct {
	Telegram      TelegramConfig `json:"telegram,omitempty"`
	ProcessList   []string       `json:"processList"`
	Logger        bool           `json:"logger,omitempty"`
	Interval      time.Duration  `json:"interval,omitempty"`
	NotifyAtStart bool           `json:"notifyAtStart,omitempty"`
	filename      string
}

const (
	// ConfigFileName is the name of config file
	ConfigFileName = "config.json"
)

var (
	config = flag.String("config", ConfigFileName, "Config file")
)

// Load configuration
func Load(config *string) (*ConfigFile, error) {
	configFile := ConfigFile{
		filename: filepath.Join("./", *config),
	}

	if _, err := os.Stat(configFile.filename); err == nil {
		file, err := os.Open(configFile.filename)
		if err != nil {
			return &configFile, err
		}
		defer file.Close()
		err = configFile.LoadFromReader(file)

		return &configFile, err
	} else if !os.IsNotExist(err) {
		return &configFile, err
	}

	return &configFile, fmt.Errorf("Config file not found")
}

// LoadFromReader yep
func (configFile *ConfigFile) LoadFromReader(configData io.Reader) error {
	if err := json.NewDecoder(configData).Decode(&configFile); err != nil {
		return err
	}

	return nil
}
