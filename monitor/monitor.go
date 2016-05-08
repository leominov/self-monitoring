package monitor

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/leominov/self-monitoring/config"
	"github.com/leominov/self-monitoring/gomonversion"
	"github.com/leominov/self-monitoring/msignal"
	"gopkg.in/telegram-bot-api.v4"
)

var (
	msgText, msgType string
	// Gomon instance
	Gomon = Monitor{}
)

const (
	// StateON for service
	StateON = "ON"
	// StateOFF for service
	StateOFF = "OFF"
	msgMask  = "%s%s switch status to %s"
)

// Service structure
type Service struct {
	Name                   string
	CurrentState, NewState bool
	DateWatch, DateUpdate  int32
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

// PrepareServiceList for list
func (monitor *Monitor) PrepareServiceList() {
	serviceList := []Service{}

	for _, name := range monitor.Config.ProcessList {
		serviceList = append(serviceList, Service{
			name,
			true,
			false,
			int32(time.Now().Unix()),
			0,
		})
	}

	monitor.ServiceList = serviceList
}

// GetPrefix for messages
func (monitor *Monitor) GetPrefix() string {
	prefixName := ""
	if monitor.Config.NodeName != "" {
		prefixName = monitor.Config.NodeName + ": "
	}

	return prefixName
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
		msgType = StateON
	}

	if len(monitor.ListOff) > 0 {
		msgText = strings.Join(append(monitor.ListOff), ", ")
		msgType = StateOFF
	}

	if msgText != "" {
		logrus.Infof(msgMask, monitor.GetPrefix(), msgText, msgType)
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

	bot, err := tgbotapi.NewBotAPI(telegram.Token)

	if err != nil {
		return err
	}

	if logrus.GetLevel() == logrus.DebugLevel {
		bot.Debug = true
	}

	if len(monitor.ListOn) > 0 {
		msgText = strings.Join(append(monitor.ListOn), ", ")
		msgType = StateON
	}

	if len(monitor.ListOff) > 0 {
		msgText = strings.Join(append(monitor.ListOff), ", ")
		msgType = StateOFF
	}

	if msgText != "" {
		msg := tgbotapi.NewMessage(telegram.ContactID, fmt.Sprintf(msgMask, monitor.GetPrefix(), msgText, msgType))

		if _, err := bot.Send(msg); err != nil {
			return fmt.Errorf("Error sending message: %v", err)
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

			if monitor.Counter > 1 {
				monitor.ServiceList[ID].DateUpdate = int32(time.Now().Unix())
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
	msgText = ""
	monitor.ListOn = []string{}
	monitor.ListOff = []string{}
}

// Configure monitor
func (monitor *Monitor) Configure() {
	config, err := config.Load(config.FileFlag)

	if err != nil {
		logrus.Errorf("Error configuring application: %v", err)
		return
	}

	config.ParseLoggerFlags()

	monitor.Config = config
	monitor.PrepareServiceList()
}

// MonitorRoutine loop
func (monitor *Monitor) MonitorRoutine() {
	var err error
	var interval time.Duration

	for {
		monitor.Counter++
		err = monitor.UpdateServiceList()

		if err != nil {
			logrus.Info(err)
			continue
		}

		monitor.CheckStatusList()

		monitor.Switch()
		monitor.Notify()

		monitor.EmptyTemp()

		interval, err = time.ParseDuration(monitor.Config.Interval)

		if err != nil {
			logrus.Infof("Error setting interval: %v", err)
			interval, _ = time.ParseDuration("5s")
		}

		time.Sleep(interval)
	}
}

// SignalRoutine loop
func (monitor *Monitor) SignalRoutine() {
	for {
		s := <-msignal.SignalChan
		switch s {
		case msignal.ReloadSignal:
			logrus.Infoln("Reloading configuration...")
			monitor.Configure()
			logrus.Infoln("Done.")

		case msignal.QuitSignal:
			logrus.Infoln("Received shutdown signal")
			msignal.ExitChan <- 0

		case msignal.InfoSignal:
			logrus.Infoln("Counter:", monitor.Counter)
			logrus.Infoln("Service list:", monitor.ServiceList)

		default:
			logrus.Infoln("Catched unknown signal")
		}
	}
}

// ExecAndNotice Execute command and Notice
func ExecAndNotice(bot *tgbotapi.BotAPI, chatID int64, command string) {
	var waitGroup sync.WaitGroup

	waitGroup.Add(1)

	go func() {
		out, err := ExecCommand(command)
		message := out
		if err != nil {
			logrus.Errorf("ExecCommand: %+v", err)
			message = err.Error()
		}

		chunks := SplitByChunk(message, 4000)
		for _, chunk := range chunks {
			if len(chunks) > 1 {
				bot.Send(tgbotapi.NewChatAction(chatID, tgbotapi.ChatTyping))
				time.Sleep(200 * time.Millisecond)
			}

			if chunk == "" {
				continue
			}

			msg := tgbotapi.NewMessage(chatID, chunk)
			bot.Send(msg)
		}
		waitGroup.Done()
	}()

	err := TimeoutWait(&waitGroup)
	if err != nil {
		msg := tgbotapi.NewMessage(chatID, err.Error())
		bot.Send(msg)
	}
}

// Control Gomon by Telegram messages
func (monitor *Monitor) Control() error {
	telegram := &monitor.Config.Telegram

	if telegram.Token == "" || telegram.ContactID == 0 {
		return fmt.Errorf("Error Telegram configuration")
	}

	bot, err := tgbotapi.NewBotAPI(telegram.Token)
	if err != nil {
		return err
	}

	if logrus.GetLevel() == logrus.DebugLevel {
		bot.Debug = true
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		return err
	}

	for update := range updates {
		isAdmin := false
		for _, name := range telegram.AdminList {
			if update.Message.From.UserName == name {
				isAdmin = true
				break
			}
		}

		if !isAdmin || !update.Message.IsCommand() {
			logrus.Errorf(
				"Access denied to exec '%s' from %s (%s)",
				update.Message.Text,
				update.Message.From.UserName,
				update.Message.Chat.Type,
			)
			continue
		}

		command := update.Message.Command()
		commandArgs := update.Message.CommandArguments()
		chatID := update.Message.Chat.ID

		/*
			sh - Exec shell command
			service - Alias for /sh service
			calc - Calculator
			up - Server uptime
			status - Service list
			who - Who is logged in
			vote - Random vote
		*/
		switch command {
		case "sh", "bash", "shell", "exec", "run":
			ExecAndNotice(bot, chatID, commandArgs)
		case "service", "srv":
			ExecAndNotice(bot, chatID, fmt.Sprintf("%s %s", "service", commandArgs))
		case "bc", "calc":
			ExecAndNotice(bot, chatID, fmt.Sprintf("echo '%s' | bc", commandArgs))
		case "w", "who":
			ExecAndNotice(bot, chatID, "who")
		case "up", "uptime":
			ExecAndNotice(bot, chatID, "uptime")
		case "st", "status":
			pref := ""
			status := ""
			for _, service := range monitor.ServiceList {
				if commandArgs != "" && commandArgs != service.Name {
					continue
				}

				state := StateON
				if service.CurrentState == false {
					state = StateOFF
				}

				status += pref + fmt.Sprintf("%s is %s", service.Name, state)
				pref = "\n"
			}
			bot.Send(tgbotapi.NewMessage(chatID, status))
		case "v", "vote":
			bot.Send(tgbotapi.NewMessage(chatID, GetVote()))
		}
	}

	return nil
}

// Run monitor
func (monitor *Monitor) Run() {
	catched, err := msignal.CatchSender()

	if err != nil {
		logrus.Fatal(err)
	}

	if catched {
		logrus.Info("Sended")
		os.Exit(0)
	}

	logrus.Debug("Debug mode enabled")

	logrus.Debugf("Starting Gomon %s...", gomonversion.Version)
	logrus.Debugf("Rinning with PID: %d", os.Getpid())

	go monitor.MonitorRoutine()
	go monitor.SignalRoutine()
	go monitor.Control()

	code := <-msignal.ExitChan
	os.Exit(code)
}
