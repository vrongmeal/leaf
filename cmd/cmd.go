// Package cmd contains the command line application.
package cmd

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/vrongmeal/leaf/pkg/engine"
	"github.com/vrongmeal/leaf/pkg/utils"
	"github.com/vrongmeal/leaf/version"

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

	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Leaf version",
		Long:  `Leaf version details`,

		Run: func(*cobra.Command, []string) {
			fmt.Printf(`%s: %s
Version [%s]
`, version.AppName, rootCmd.Short, version.Version)
		},
	}
)

func init() {
	logrus.SetLevel(logrus.TraceLevel)
	logrus.SetFormatter(new(prefixed.TextFormatter))

	rootCmd.PersistentFlags().StringVarP(&confPath, "config", "c", utils.DefaultConfPath, "Config path for leaf configuration file")

	rootCmd.AddCommand(versionCmd)
}

// Execute starts the command line tool.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logrus.Fatalf("Cannot start leaf: %s", err.Error())
	}
}
