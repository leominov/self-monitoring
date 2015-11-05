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

	"github.com/Sirupsen/logrus"
)

// TelegramConfig structure
type TelegramConfig struct {
	Enable    bool   `json:"enable,omitempty"`
	Token     string `json:"token,omitempty"`
	ContactID int    `json:"contactID,omitempty"`
}

// File structure
type File struct {
	Telegram      TelegramConfig `json:"telegram,omitempty"`
	ProcessList   []string       `json:"processList"`
	Interval      time.Duration  `json:"interval,omitempty"`
	NotifyAtStart bool           `json:"notifyAtStart,omitempty"`
	LogLevel      string         `json:"logLevel,omitempty"`
	NodeName      string         `json:"nodeName,omitempty"`
	filename      string
}

const (
	// FileName is the name of config file
	FileName = "config.json"
)

var (
	// FileFlag is the config file from flag()
	FileFlag = flag.String("config", FileName, "Config file")
	// DebugFlag global mode
	DebugFlag         = flag.Bool("debug", false, "Enable debug mode")
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

// ParseLoggerFlags for Logrus
func (configFile *File) ParseLoggerFlags() {
	if configFile.LogLevel != "" {
		lvl, err := logrus.ParseLevel(configFile.LogLevel)
		if err != nil {
			fmt.Printf("Unable to parse logging level: %s\n", configFile.LogLevel)
			os.Exit(1)
		}
		logrus.SetLevel(lvl)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}

	if *DebugFlag {
		os.Setenv("DEBUG", "1")
		logrus.SetLevel(logrus.DebugLevel)
	}
}

// LoadFromReader json
func (configFile *File) LoadFromReader(configData io.Reader) error {
	if err := json.NewDecoder(configData).Decode(&configFile); err != nil {
		return err
	}

	return nil
}
