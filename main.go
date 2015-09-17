package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/Syfaro/telegram-bot-api"
)

// Service structure
type Service struct {
	name         string
	current, new bool
}

// TelegramConfig structure
type TelegramConfig struct {
	Enable    bool   `json:"enable,omitempty"`
	Token     string `json:"token,omitempty"`
	ContactID int    `json:"contactID,omitempty"`
	Debug     bool   `json:"debug,omitempty"`
}

// ConfigFile structure
type ConfigFile struct {
	Telegram    TelegramConfig `json:"telegram,omitempty"`
	ProcessList []string       `json:"processList"`
	Logger      bool           `json:"logger,omitempty"`
	Interval    time.Duration  `json:"interval,omitempty"`
	filename    string
}

// Monitor structure
type Monitor struct {
	Config             *ConfigFile
	CurrentServiceList []string
	ServiceList        []Service
	ListOn, ListOff    []string
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

// Prepare parameters
func (monitor *Monitor) Prepare() {
	serviceList := []Service{}

	for _, name := range monitor.Config.ProcessList {
		serviceList = append(serviceList, Service{name, true, false})
	}

	monitor.ServiceList = serviceList
}

// UpdateServiceList getting current process list
func (monitor *Monitor) UpdateServiceList() error {
	cmd := exec.Command("ps", "-eo", "comm=|sort|uniq")
	output, err := cmd.CombinedOutput()

	if err != nil {
		return err
	}

	monitor.CurrentServiceList = strings.Split(string(output), "\n")

	return nil
}

// CheckStatusList for monitor
func (monitor *Monitor) CheckStatusList() []Service {
	for ID, service := range monitor.ServiceList {
		monitor.ServiceList[ID].new = false
		for _, sname := range monitor.CurrentServiceList {
			if monitor.ServiceList[ID].new == true {
				continue
			} else if sname == service.name {
				monitor.ServiceList[ID].new = true
				continue
			}
		}
	}

	return monitor.ServiceList
}

// RunLogger service status
func (monitor *Monitor) RunLogger() error {
	if len(monitor.ListOn) > 0 {
		fmt.Printf("%s: %s switch status to ON\n", time.Now(), strings.Join(append(monitor.ListOn), ", "))
	}

	if len(monitor.ListOff) > 0 {
		fmt.Printf("%s: %s switch status to OFF\n", time.Now(), strings.Join(append(monitor.ListOff), ", "))
	}

	return nil
}

// RunTelegram service status
func (monitor *Monitor) RunTelegram() error {
	telegram := &monitor.Config.Telegram

	if telegram.Token == "" || telegram.ContactID == 0 {
		fmt.Println("Error. Check configuration parameters:")
		fmt.Println(" - Telegram.Token")
		fmt.Println(" - Telegram.ContactID")
		return fmt.Errorf("Error Telegram configuration")
	}

	bot, err := tgbotapi.NewBotAPI(monitor.Config.Telegram.Token)

	if err != nil {
		return err
	}

	bot.Debug = telegram.Debug

	if len(monitor.ListOn) > 0 {
		msg := tgbotapi.NewMessage(telegram.ContactID, fmt.Sprintf("%s switch status to ON", strings.Join(append(monitor.ListOn), ", ")))
		bot.SendMessage(msg)
	}

	if len(monitor.ListOff) > 0 {
		msg := tgbotapi.NewMessage(telegram.ContactID, fmt.Sprintf("%s switch status to OFF", strings.Join(append(monitor.ListOff), ", ")))
		bot.SendMessage(msg)
	}

	return nil
}

// Switch service status
func (monitor *Monitor) Switch() {
	for ID, service := range monitor.ServiceList {
		if service.current != service.new {
			if service.new {
				monitor.ListOn = append(monitor.ListOn, service.name)
			} else {
				monitor.ListOff = append(monitor.ListOff, service.name)
			}

			monitor.ServiceList[ID].current = service.new
		}
	}
}

// Notify service status
func (monitor *Monitor) Notify() {
	if monitor.Config.Logger {
		monitor.RunLogger()
	}

	if monitor.Config.Telegram.Enable {
		monitor.RunTelegram()
	}
}

// EmptyTemp data
func (monitor *Monitor) EmptyTemp() {
	monitor.ListOn = []string{}
	monitor.ListOff = []string{}
}

// Run monitor
func (monitor *Monitor) Run() error {
	err := monitor.UpdateServiceList()

	if err != nil {
		return err
	}

	monitor.CheckStatusList()

	monitor.Switch()
	monitor.Notify()

	monitor.EmptyTemp()

	return nil
}

func main() {
	flag.Parse()
	config, err := Load(config)

	if err != nil {
		log.Panic(err)
		return
	}

	monitor := Monitor{
		config,
		[]string{},
		[]Service{},
		[]string{},
		[]string{},
	}

	monitor.Prepare()

	for {
		monitor.Run()

		time.Sleep(monitor.Config.Interval * time.Second)
	}
}
