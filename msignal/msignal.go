package msignal

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/Sirupsen/logrus"
	"github.com/leominov/self-monitoring/config"
)

const (
	// ReloadSignal for reloading configuration
	ReloadSignal = syscall.SIGUSR1
	// QuitSignal for exit
	QuitSignal = syscall.SIGQUIT
	// InfoSignal for additional information
	InfoSignal = syscall.SIGINFO
)

var (
	// SignalChan for system signals
	SignalChan = make(chan os.Signal, 1)
	// ExitChan for exit chan
	ExitChan = make(chan int)
)

// CatchSender for catching signar request
func CatchSender() {
	if *config.SignalFlag != "" {
		p, err := os.FindProcess(*config.PidFlag)
		if err != nil {
			logrus.Errorf("Error sending signal: %s", err)
		}

		switch *config.SignalFlag {
		case "reload":
			p.Signal(ReloadSignal)
		case "quit":
			p.Signal(QuitSignal)
		case "info":
			p.Signal(InfoSignal)
		}

		logrus.Info("OK")
		os.Exit(0)
	}

	signal.Notify(SignalChan,
		InfoSignal,
		ReloadSignal,
		QuitSignal)
}
