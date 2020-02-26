// Package cmd contains the command line application.
package cmd

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/vrongmeal/leaf/pkg/engine"
	"github.com/vrongmeal/leaf/pkg/types"
	"github.com/vrongmeal/leaf/pkg/utils"
	"github.com/vrongmeal/leaf/version"

	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

var (
	confPath string
	conf     types.Config

	rootCmd = &cobra.Command{
		Use:   "leaf",
		Short: "General purpose hot-reloader for all projects",
		Long: `Given a set of commands, leaf watches the filtered paths in the project directory for any changes and runs the commands in
order so you don't have to yourself`,

		PreRun: func(*cobra.Command, []string) {
			if confPath != "" {
				viper.SetConfigFile(confPath)
			} else {
				viper.SetConfigFile(utils.DefaultConfPath)
			}

			viper.AutomaticEnv()

			if err := viper.ReadInConfig(); err != nil {
				logrus.Warnln("Continuing without config...")
			}

			if err := viper.Unmarshal(&conf); err != nil {
				logrus.Fatalf("Cannot unmarshal config: %s", err.Error())
			}

			// Check and include defaults if required
			if len(conf.Exclude) == 0 {
				conf.Exclude = utils.DefaultExcludePaths
			} else {
				finalExcludes := []string{}
				for _, e := range conf.Exclude {
					if e == utils.DefaultExcludePathsKeyword {
						finalExcludes = append(finalExcludes, utils.DefaultExcludePaths...)
						continue
					}

					finalExcludes = append(finalExcludes, e)
				}

				conf.Exclude = finalExcludes
			}
		},

		Run: func(*cobra.Command, []string) {
			logrus.Infof("Starting to watch: %s", conf.Root)
			logrus.Infoln("Excluded paths:")
			for _, e := range conf.Exclude {
				logrus.Infof("\t%s", e)
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
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(new(prefixed.TextFormatter))

	rootCmd.PersistentFlags().StringVarP(&confPath, "config", "c", utils.DefaultConfPath, "Config path for leaf configuration file")

	rootCmd.Flags().StringP("root", "r", utils.CWD, "Root directory to watch")
	rootCmd.Flags().StringSliceP(
		"exclude", "e", utils.DefaultExcludePaths,
		fmt.Sprintf("Paths to exclude. You can append default paths by adding `%s` in your list", utils.DefaultExcludePathsKeyword))
	rootCmd.Flags().StringSliceP("filters", "f", []string{}, "Filters to apply to watch")
	rootCmd.Flags().StringSliceP("exec", "x", []string{}, "Exec commands on file change")
	rootCmd.Flags().DurationP("delay", "d", 500*time.Millisecond, "Delay after which commands are run on file change")

	if err := viper.BindPFlag("root", rootCmd.Flags().Lookup("root")); err != nil {
		logrus.Fatalf("Error binding flag to viper: %s", err.Error())
	}
	if err := viper.BindPFlag("exclude", rootCmd.Flags().Lookup("exclude")); err != nil {
		logrus.Fatalf("Error binding flag to viper: %s", err.Error())
	}
	if err := viper.BindPFlag("filters", rootCmd.Flags().Lookup("filters")); err != nil {
		logrus.Fatalf("Error binding flag to viper: %s", err.Error())
	}
	if err := viper.BindPFlag("exec", rootCmd.Flags().Lookup("exec")); err != nil {
		logrus.Fatalf("Error binding flag to viper: %s", err.Error())
	}
	if err := viper.BindPFlag("delay", rootCmd.Flags().Lookup("delay")); err != nil {
		logrus.Fatalf("Error binding flag to viper: %s", err.Error())
	}

	rootCmd.AddCommand(versionCmd)
}

// Execute starts the command line tool.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logrus.Fatalf("Cannot start leaf: %s", err.Error())
	}
}
