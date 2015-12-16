package main

import (
	"flag"

	"github.com/leominov/self-monitoring/monitor"
	"github.com/leominov/self-monitoring/msignal"
)

func main() {
	flag.Parse()
	msignal.CatchSender()

	monitor := monitor.Monitor{}

	monitor.Configure()
	monitor.Run()
}
