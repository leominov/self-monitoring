// +build !windows

package monitor

import (
	"github.com/Sirupsen/logrus"
	"github.com/leominov/self-monitoring/msignal"
)

// SignalRoutine loop
func (monitor *Monitor) SignalRoutine() {
	for {
		s := <-msignal.SignalChan
		switch s {
		case msignal.ReloadSignal:
			logrus.Infoln("Reloading configuration...")
			monitor.Configure()
			logrus.Infoln("Done.")

		case msignal.QuitSignal:
			logrus.Infoln("Received shutdown signal")
			msignal.ExitChan <- 0

		case msignal.InfoSignal:
			logrus.Infoln("Counter:", monitor.Counter)
			logrus.Infoln("Service list:", monitor.ServiceList)

		default:
			logrus.Infoln("Catched unknown signal")
		}
	}
}
