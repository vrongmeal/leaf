package engine

import (
	"os"
	"os/signal"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/vrongmeal/leaf/pkg/commander"
	"github.com/vrongmeal/leaf/pkg/utils"
	"github.com/vrongmeal/leaf/pkg/watcher"
)

// Start runs the watcher and executes the commands from the config on file change.
func Start(conf *utils.Config) error {
	cmdr := commander.NewCommander(conf.Exec)

	opts := watcher.WatchOpts{
		Root:    conf.Root,
		Exclude: conf.Exclude,
		Filters: conf.Filters,
	}

	wchr, err := watcher.NewWatcher(&opts)
	if err != nil {
		return err
	}

	go func() {
		if err := wchr.Watch(); err != nil {
			logrus.Errorf("Failed to setup watcher: %s", err.Error())
		}
	}()

	interrupt := make(chan os.Signal)
	signal.Notify(interrupt, os.Interrupt)

	exit := make(chan bool, 1)

	runCmd(cmdr)

	go func() {
		<-interrupt
		logrus.Infoln("Terminating after cleanup")
		if err := cmdr.Kill(); err != nil {
			logrus.Errorf("Error while stopping command: %s", err.Error())
		}
		wchr.Close()
		exit <- true
	}()

	for {
		select {
		case file := <-wchr.File:
			if err := cmdr.Kill(); err != nil {
				logrus.Errorf("Error while stopping command: %s", err.Error())
				wchr.Close()
				return err
			}
			logrus.Infof("File modified! Reloading... (%s)", file)

			// Sleep for conf.Delay duration amount of time.
			time.Sleep(conf.Delay)

			runCmd(cmdr)

		case err := <-wchr.Err:
			logrus.Errorf("Error while watching: %s", err.Error())

		case <-exit:
			return nil
		}
	}
}

func runCmd(cmdr *commander.Commander) {
	go func() {
		if err := cmdr.Run(); err != nil {
			logrus.Warnf("Error while running command: %s", err.Error())
		}
	}()
}
