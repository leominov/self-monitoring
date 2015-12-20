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
	errSign  error
)

// CatchSender for catching signar request
func CatchSender() (bool, error) {
	if *config.SignalFlag != "" {
		p, err := os.FindProcess(*config.PidFlag)

		if err != nil {
			return false, err
		}

		if p.Pid == 0 {
			err := fmt.Errorf("Process with pid %d not found", *config.PidFlag)
			return false, err
		}

		switch *config.SignalFlag {
		case "reload":
			errSign = p.Signal(ReloadSignal)
		case "quit":
			errSign = p.Signal(QuitSignal)
		case "info":
			errSign = p.Signal(InfoSignal)
		default:
			return false, errors.New("Unknown signal")
		}

		if errSign != nil {
			err := fmt.Errorf("Error sending signal: %v", errSign)
			return false, err
		}

		return true, nil
	}

	signal.Notify(SignalChan,
		InfoSignal,
		ReloadSignal,
		QuitSignal)

	return false, nil
}
