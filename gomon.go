package main

import (
	"flag"
	"log"

	"github.com/leominov/self-monitoring/config"
	"github.com/leominov/self-monitoring/monitor"
)

func main() {
	flag.Parse()
	config, err := config.Load(config.FileFlag)

	if err != nil {
		log.Panic(err)
		return
	}

	monitor := monitor.Monitor{}
	monitor.Config = config

	monitor.Run()
}
