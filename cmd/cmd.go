// Package cmd implements the command-line interface for the
// leaf command. It contains commands and their flags.
package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/vrongmeal/leaf"

	lpf "github.com/x-cray/logrus-prefixed-formatter"
)

var (
	confPath string
	debugEnv bool

	conf leaf.Config
)

var rootCmd = &cobra.Command{
	Use:   "leaf",
	Short: "general purpose hot-reloader for all projects",
	Long: `
Command leaf watches for changes in the working directory and
runs the specified set of commands whenever a file updates.
A set of filters can be applied to the watch and directories
can be excluded.`,

	PersistentPreRun: func(*cobra.Command, []string) {
		// Logger is initialized depending upon the debug
		// flag and hence is not in the `init` function.
		initialiseLogger()
	},

	PreRun: func(*cobra.Command, []string) {
		ferr, rerr := setupConfig()
		if rerr != nil {
			log.Fatalln(rerr)
		} else if ferr != nil {
			log.Warnf("config file not read: %v", ferr)
		}
	},

	Run: func(*cobra.Command, []string) {
		log.Infof("watching '%s'", conf.Root)

		if err := runEngine(&conf); err != nil {
			log.Fatalln(err)
		}
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "prints leaf version",
	Long: `
Prints the build information for the leaf commnd line.`,

	Run: func(*cobra.Command, []string) {
		goModInfo, err := leaf.GoModuleInfo()
		if err != nil {
			log.Fatalf("error getting version: %v", err)
		}

		fmt.Printf("leaf version %s\n", goModInfo.Version)
	},
}

func init() {
	initializeFlags()

	if err := bindFlagsToConfig(); err != nil {
		log.Fatalf("cannot bind flags with config: %v", err)
	}

	rootCmd.AddCommand(versionCmd)
}

// Execute starts the command line tool.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}

// initialiseLogger sets up the logger configuration.
func initialiseLogger() {
	if debugEnv {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	log.SetFormatter(&lpf.TextFormatter{
		DisableTimestamp: true,
	})
}

// initializeFlags sets the initializes flags for the commands.
func initializeFlags() {
	rootCmd.PersistentFlags().BoolVar(
		&debugEnv, "debug", false,
		"run in development (debug) environment")

	rootCmd.PersistentFlags().StringVarP(
		&confPath, "config", "c", leaf.DefaultConfPath,
		"config path for the configuration file")

	rootCmd.Flags().StringP(
		"root", "r", leaf.CWD,
		"root directory to watch")

	rootCmd.Flags().StringSliceP(
		"exclude", "e", leaf.DefaultExcludePaths,
		"paths to exclude from watching")

	rootCmd.Flags().StringSliceP(
		"filters", "f", []string{},
		"filters to apply to watch")

	rootCmd.Flags().StringSliceP(
		"exec", "x", []string{},
		"exec commands on file change")

	rootCmd.Flags().DurationP(
		"delay", "d", 500*time.Millisecond,
		"delay after which commands are run on file change")
}

// bindFlagsToConfig binds the flags with viper config file.
func bindFlagsToConfig() error {
	keyFlagMap := map[string]string{
		"root":    "root",
		"exclude": "exclude",
		"filters": "filters",
		"exec":    "exec",
		"delay":   "delay",
	}

	for key, flag := range keyFlagMap {
		err := viper.BindPFlag(key, rootCmd.Flags().Lookup(flag))
		if err != nil {
			return err
		}
	}

	return nil
}

// setupConfig reads and unmarshals the config file into
// the `conf` variable.
func setupConfig() (fileErr, readErr error) {
	if confPath != "" {
		viper.SetConfigFile(confPath)
	} else {
		viper.SetConfigFile(leaf.DefaultConfPath)
	}

	viper.AutomaticEnv()

	var confFileErr error
	if err := viper.ReadInConfig(); err != nil {
		// Even if no config is provided we still unmarshal
		// the config because are flags are bound with conf.
		confFileErr = err
	}

	if err := viper.Unmarshal(&conf); err != nil {
		return confFileErr, fmt.Errorf("unable to read config: %v", err)
	}

	// By default the defaults are included in the excludes.
	// If the excludes are specified explicitly, defaults will
	// only be included if the `DEFAULTS` keyword is in the
	// excluded paths.
	if len(conf.Exclude) == 0 {
		conf.Exclude = leaf.DefaultExcludePaths
	} else {
		finalExcludes := []string{}
		for _, e := range conf.Exclude {
			if e == leaf.DefaultExcludePathsKeyword {
				finalExcludes = append(finalExcludes,
					leaf.DefaultExcludePaths...)
				continue
			}

			finalExcludes = append(finalExcludes, e)
		}

		conf.Exclude = finalExcludes
	}

	return confFileErr, nil
}

// newCmdContext returns a context which cancels on an OS
// interrupt, i.e., cancels when process is killed.
func newCmdContext(onInterrupt func(os.Signal)) context.Context {
	interrupt := make(chan os.Signal)
	signal.Notify(interrupt, os.Interrupt)
	ctx, cancel := context.WithCancel(context.Background())

	go func(onSignal func(os.Signal), cancelProcess context.CancelFunc) {
		sig := <-interrupt
		onSignal(sig)
		cancelProcess()
	}(onInterrupt, cancel)

	return ctx
}

// runEngine runs the watcher and executes the commands from
// the config on file change.
func runEngine(conf *leaf.Config) error {
	ctx := newCmdContext(func(s os.Signal) {
		log.Infof("closing: signal received: %s", s.String())
	})

	commander := leaf.NewCommander(leaf.Commander{
		Commands: conf.Exec,
		OnStart: func(cmd *leaf.Command) {
			log.Infof("running: %s", cmd.String())
		},
		OnError: func(err error) {
			log.Errorln(err)
		},
		OnExit: func() {
			log.Info("commands executed")
		},
		ExitOnError: false,
	})

	watcher, err := leaf.NewWatcher(
		conf.Root, conf.Exclude, conf.Filters)
	if err != nil {
		log.Fatalf("error creating watcher: %v", err)
	}

	cmdCtx, killCmds := context.WithCancel(ctx)
	go commander.Run(cmdCtx)

	for wr := range watcher.Watch(ctx) {
		if wr.Err != nil {
			log.Errorf("error while watching: %v", err)
			continue
		}

		log.Infof("file '%s' changed, reloading...", wr.File)

		killCmds()                                 // kill previous commands
		cmdCtx, killCmds = context.WithCancel(ctx) // new context
		time.Sleep(conf.Delay)                     // wait for 'delay' duration
		<-commander.Done()                         // wait more if required by commands
		go commander.Run(cmdCtx)                   // run commands
	}

	killCmds()
	<-commander.Done()

	log.Infoln("shutdown successfully")
	return nil
}
