package monitor

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/Syfaro/telegram-bot-api"
	"github.com/leominov/self-monitoring/config"
)

var (
	msgText, msgType string
)

// Service structure
type Service struct {
	Name                   string
	CurrentState, NewState bool
}

// Monitor structure
type Monitor struct {
	Config             *config.File
	CurrentServiceList []string
	ServiceList        []Service
	ListOn, ListOff    []string
	Counter            int
	// Telegram           *tgbotapi.BotAPI
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
		monitor.ServiceList[ID].NewState = false
		for _, sname := range monitor.CurrentServiceList {
			if monitor.ServiceList[ID].NewState == true {
				continue
			} else if sname == service.Name {
				monitor.ServiceList[ID].NewState = true
				continue
			}
		}
	}

	return monitor.ServiceList
}

// RunLogger service status
func (monitor *Monitor) RunLogger() error {
	if len(monitor.ListOn) > 0 {
		msgText = strings.Join(append(monitor.ListOn), ", ")
		msgType = "ON"
	}

	if len(monitor.ListOff) > 0 {
		msgText = strings.Join(append(monitor.ListOff), ", ")
		msgType = "ON"
	}

	if msgText != "" {
		logrus.Infof("%s switch status to %s", msgText, msgType)
	}

	return nil
}

// RunTelegram service status
func (monitor *Monitor) RunTelegram() error {
	telegram := &monitor.Config.Telegram

	if telegram.Token == "" || telegram.ContactID == 0 {
		logrus.Debug("Check Telegram config parameters: token, contactID")
		return fmt.Errorf("Error Telegram configuration")
	}

	bot, err := tgbotapi.NewBotAPI(monitor.Config.Telegram.Token)

	if err != nil {
		return err
	}

	if logrus.GetLevel() == logrus.DebugLevel {
		bot.Debug = true
	}

	if len(monitor.ListOn) > 0 {
		msgText = strings.Join(append(monitor.ListOn), ", ")
		msgType = "ON"
	}

	if len(monitor.ListOff) > 0 {
		msgText = strings.Join(append(monitor.ListOff), ", ")
		msgType = "OFF"
	}

	if msgText != "" {
		msg := tgbotapi.NewMessage(telegram.ContactID, fmt.Sprintf("%s switch status to %s", msgText, msgType))

		if _, err := bot.SendMessage(msg); err != nil {
			return fmt.Errorf("Error sending message: %s", err)
		}
	}

	return nil
}

// Switch service status
func (monitor *Monitor) Switch() {
	for ID, service := range monitor.ServiceList {
		if service.CurrentState != service.NewState {
			if service.NewState {
				monitor.ListOn = append(monitor.ListOn, service.Name)
			} else {
				monitor.ListOff = append(monitor.ListOff, service.Name)
			}

			monitor.ServiceList[ID].CurrentState = service.NewState
		}
	}
}

// Notify service status
func (monitor *Monitor) Notify() {
	if monitor.Config.NotifyAtStart == false && monitor.Counter == 1 {
		return
	}

	monitor.RunLogger()

	if monitor.Config.Telegram.Enable {
		if err := monitor.RunTelegram(); err != nil {
			logrus.Error(err)
		}
	}
}

// EmptyTemp data
func (monitor *Monitor) EmptyTemp() {
	monitor.ListOn = []string{}
	monitor.ListOff = []string{}
}

// Run monitor
func (monitor *Monitor) Run() {
	logrus.Debug("Starting Gomon...")

	monitor.Prepare()

	for {
		monitor.Counter++
		err := monitor.UpdateServiceList()

		if err != nil {
			log.Print(err)
			continue
		}

		monitor.CheckStatusList()

		monitor.Switch()
		monitor.Notify()

		monitor.EmptyTemp()

		time.Sleep(monitor.Config.Interval * time.Millisecond)
	}
}
