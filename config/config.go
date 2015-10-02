package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// TelegramConfig structure
type TelegramConfig struct {
	Enable    bool   `json:"enable,omitempty"`
	Token     string `json:"token,omitempty"`
	ContactID int    `json:"contactID,omitempty"`
	Debug     bool   `json:"debug,omitempty"`
}

// File structure
type File struct {
	Telegram      TelegramConfig `json:"telegram,omitempty"`
	ProcessList   []string       `json:"processList"`
	Logger        bool           `json:"logger,omitempty"`
	Interval      time.Duration  `json:"interval,omitempty"`
	NotifyAtStart bool           `json:"notifyAtStart,omitempty"`
	filename      string
}

const (
	// FileName is the name of config file
	FileName = "config.json"
)

var (
	// FileFlag is the config file from flag()
	FileFlag          = flag.String("config", FileName, "Config file")
	loadWD, loadWDErr = os.Getwd()
)

// Load configuration
func Load(config *string) (*File, error) {
	if loadWDErr == nil {
		loadWD = strings.Replace(loadWD, "\\", "/", -1) + "/"
		if filepath.Dir(*config) == "." {
			*config = loadWD + *config
		}
	}

	configFile := File{
		filename: *config,
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
func (configFile *File) LoadFromReader(configData io.Reader) error {
	if err := json.NewDecoder(configData).Decode(&configFile); err != nil {
		return err
	}

	return nil
}
