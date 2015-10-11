package main

import (
	"flag"

	"github.com/Sirupsen/logrus"
	"github.com/leominov/self-monitoring/config"
	"github.com/leominov/self-monitoring/monitor"
)

func main() {
	flag.Parse()
	config, err := config.Load(config.FileFlag)
	config.ParseLoggerFlags()

	if err != nil {
		logrus.Errorf("Error configuring application: %s", err)
		return
	}

	monitor := monitor.Monitor{}
	monitor.Config = config

	monitor.Run()
}
