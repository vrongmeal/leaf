package main

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/vrongmeal/leaf/pkg/engine"
	"github.com/vrongmeal/leaf/pkg/utils"

	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

var (
	confPath string

	rootCmd = &cobra.Command{
		Use:   "leaf",
		Short: "General purpose hot-reloader for all projects",
		Long: `Given a set of commands, leaf watches the filtered paths in the project directory for any changes and runs the commands in
order so you don't have to yourself`,

		Run: func(*cobra.Command, []string) {
			conf, err := utils.GetConfig(confPath)
			if err != nil {
				logrus.Fatalf("Error getting config: %s", err.Error())
			}

			isdir, err := utils.IsDir(conf.Root)
			if err != nil || !isdir {
				conf.Root = utils.CWD
			}

			logrus.Infof("Starting to watch: %s", conf.Root)
			logrus.Infoln("Excluded paths:")
			for i, e := range conf.Exclude {
				logrus.Infof("%d. %s", i, e)
			}

			if err := engine.Start(&conf); err != nil {
				logrus.Fatalf("Cannot start leaf: %s", err.Error())
			}
		},
	}
)

func init() {
	logrus.SetLevel(logrus.TraceLevel)
	logrus.SetFormatter(new(prefixed.TextFormatter))

	rootCmd.PersistentFlags().StringVarP(&confPath, "config", "c", utils.DefaultConfPath, "Config path for leaf configuration file")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		logrus.Fatalf("Cannot start leaf: %s", err.Error())
	}
}
