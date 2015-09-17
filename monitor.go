package main

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/Syfaro/telegram-bot-api"
)

// Service structure
type Service struct {
	name         string
	current, new bool
}

// Monitor structure
type Monitor struct {
	Config             *ConfigFile
	CurrentServiceList []string
	ServiceList        []Service
	ListOn, ListOff    []string
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
func (monitor *Monitor) Run() {
	for {
		err := monitor.UpdateServiceList()

		if err != nil {
			fmt.Println(err)
			continue
		}

		monitor.CheckStatusList()

		monitor.Switch()
		monitor.Notify()

		monitor.EmptyTemp()

		time.Sleep(monitor.Config.Interval * time.Second)
	}
}
