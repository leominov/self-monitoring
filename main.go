package main

import (
	"flag"
	"log"
)

func main() {
	flag.Parse()
	config, err := Load(config)

	if err != nil {
		log.Panic(err)
		return
	}

	monitor := Monitor{}
	monitor.Config = config

	monitor.Run()
}
