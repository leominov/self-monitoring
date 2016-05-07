package monitor

import (
	"errors"
	"os/exec"
	"sync"
	"time"
)

var (
	// ErrShellEmpty for Control
	ErrShellEmpty = errors.New("Command not found.")
	// ErrShellTimeout for Control
	ErrShellTimeout = errors.New("Timeout")
)

// TimeoutWait for ExecCommand
func TimeoutWait(waitGroup *sync.WaitGroup) error {
	// Make a chanel that we'll use as a timeout
	c := make(chan int, 1)

	// Start waiting for the routines to finish
	go func() {
		waitGroup.Wait()
		c <- 1
	}()

	select {
	case _ = <-c:
		return nil
	case <-time.After(10 * time.Second):
		return ErrShellTimeout
	}
}

// ExecCommand for control
func ExecCommand(cmd string) (string, error) {
	if cmd == "" {
		return "", ErrShellEmpty
	}

	out, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		return "", err
	}

	return string(out), nil
}
