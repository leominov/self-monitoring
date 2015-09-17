package main

import (
	"flag"
	"log"
	"time"
)

func main() {
	flag.Parse()
	config, err := Load(config)

	if err != nil {
		log.Panic(err)
		return
	}

	monitor := Monitor{
		config,      // Config
		[]string{},  // CurrentServiceList
		[]Service{}, // ServiceList
		[]string{},  // ListOn
		[]string{},  // ListOff
	}

	monitor.Prepare()

	for {
		monitor.Run()
		time.Sleep(monitor.Config.Interval * time.Second)
	}
}
