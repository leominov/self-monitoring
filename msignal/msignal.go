package msignal

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

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
func CatchSender() (bool, error) {
	if *config.SignalFlag != "" {
		var err error
		var p *os.Process

		p, err = os.FindProcess(*config.PidFlag)

		if err != nil {
			return false, err
		}

		if p.Pid == 0 {
			return false, fmt.Errorf("Process with pid %d not found", *config.PidFlag)
		}

		switch *config.SignalFlag {
		case "reload":
			err = p.Signal(ReloadSignal)
		case "quit":
			err = p.Signal(QuitSignal)
		case "info":
			err = p.Signal(InfoSignal)
		default:
			return false, errors.New("Unknown signal")
		}

		if err != nil {
			return false, fmt.Errorf("Error sending signal: %v", err)
		}

		return true, nil
	}

	signal.Notify(SignalChan,
		InfoSignal,
		ReloadSignal,
		QuitSignal)

	return false, nil
}
