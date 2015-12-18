package main

import (
	"flag"

	"github.com/leominov/self-monitoring/monitor"
)

func main() {
	flag.Parse()

	monitor.Gomon.Configure()
	monitor.Gomon.Run()
}
