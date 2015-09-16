package main

import (
	"flag"
	"fmt"
	"log"
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

var (
	processList = []string{
		"acrypt",
		"capella",
		"docker",
		"igmpproxy",
		"memcached",
		"mongod",
		"mysqld",
		"nginx",
		"openvpn",
		"php5-fpm",
		"postgres",
		"pptpd",
		"redis-server",
		"smbd",
		"transmission",
		"vsftpd",
		"top",
	}

	token     = flag.String("token", "", "Telegram API Token")
	contactID = flag.Int("client", 0, "Telegram Contact ID")
)

func createList() []Service {
	serviceList := []Service{}

	for _, name := range processList {
		serviceList = append(serviceList, Service{name, true, false})
	}

	return serviceList
}

func getServiceList() ([]string, error) {
	cmd := exec.Command("ps", "-eo", "comm=|sort|uniq")
	output, err := cmd.CombinedOutput()

	if err != nil {
		return nil, err
	}

	return strings.Split(string(output), "\n"), nil
}

func checkStatusList(srvList []string, serviceList []Service) []Service {
	for ID, service := range serviceList {
		serviceList[ID].new = false
		for _, sname := range srvList {
			if serviceList[ID].new == true {
				continue
			} else if sname == service.name {
				serviceList[ID].new = true
				continue
			}
		}
	}

	return serviceList
}

func logSwitch(listOn []string, listOff []string) {
	if len(listOn) > 0 {
		fmt.Printf("%s: %s switch status to ON\n", time.Now(), strings.Join(append(listOn), ", "))
	}

	if len(listOff) > 0 {
		fmt.Printf("%s: %s switch status to OFF\n", time.Now(), strings.Join(append(listOff), ", "))
	}
}

func notifySwitch(listOn []string, listOff []string, bot *tgbotapi.BotAPI) {
	if len(listOn) > 0 {
		msg := tgbotapi.NewMessage(*contactID, fmt.Sprintf("%s switch status to ON", strings.Join(append(listOn), ", ")))
		bot.SendMessage(msg)
	}

	if len(listOff) > 0 {
		msg := tgbotapi.NewMessage(*contactID, fmt.Sprintf("%s switch status to OFF", strings.Join(append(listOff), ", ")))
		bot.SendMessage(msg)
	}
}

func switchAndNotify(serviceList []Service, bot *tgbotapi.BotAPI) {
	listOn := []string{}
	listOff := []string{}

	for ID, service := range serviceList {
		if service.current != service.new {
			if service.new {
				listOn = append(listOn, service.name)
			} else {
				listOff = append(listOff, service.name)
			}

			serviceList[ID].current = service.new
		}
	}

	logSwitch(listOn, listOff)
	notifySwitch(listOn, listOff, bot)
}

func main() {
	flag.Parse()

	if *token == "" || *contactID == 0 {
		flag.PrintDefaults()
		return
	}

	serviceList := createList()
	bot, err := tgbotapi.NewBotAPI(*token)

	if err != nil {
		log.Panic(err)
	}

	for {
		currentServiceList, err := getServiceList()

		if err == nil {
			srv := checkStatusList(currentServiceList, serviceList)
			switchAndNotify(srv, bot)
		}

		time.Sleep(3 * time.Second)
	}
}
