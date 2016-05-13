// +build windows

package msignal

import "os"

var (
	// SignalChan for system signals
	SignalChan = make(chan os.Signal, 1)
	// ExitChan for exit chan
	ExitChan = make(chan int)
)

// CatchSender is skipped for Windows
func CatchSender() (bool, error) {
	return false, nil
}
