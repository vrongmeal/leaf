// Package leaf provides with utilities to create the leaf
// CLI tool. It includes watcher, filters and commander which
// watch files for changes, filter out required results and
// execute external commands respectively.
//
// The package comes with utilities that can aid in creating
// a reloader with a simple go program.
//
// Let's look at an example where the watcher watches the `src/`
// directory for changes and for any changes builds the project.
//
// 	package main
//
// 	import (
// 		"log"
// 		"os"
// 		"path/filepath"
//
// 		"github.com/vrongmeal/leaf"
// 	)
//
// 	func main() {
// 		// Use a context that cancels when program is interrupted.
// 		ctx := leaf.NewCmdContext(func(os.Signal) {
// 			log.Println("Shutting down.")
// 		})
//
// 		cwd, err := os.Getwd()
// 		if err != nil {
// 			log.Fatalln(err)
// 		}
//
// 		// Root is <cwd>/src
// 		root := filepath.Join(cwd, "src")
//
// 		// Exclude "src/index.html" from results.
// 		filters := []leaf.Filter{
// 			{Include: false, Pattern: "src/index.html"},
// 		}
//
// 		filterCollection := leaf.NewFilterCollection(
// 			filters,
// 			// Matches directory or filepath.Match expressions
// 			leaf.StandardFilterMatcher,
// 			// Definitely excludes and shows only includes (if any)
// 			leaf.StandardFilterHandler)
//
// 		watcher, err := leaf.NewWatcher(
// 			root,
// 			// Standard paths to exclude, like vendor, .git,
// 			// node_modules, venv etc.
// 			leaf.DefaultExcludePaths,
// 			filterCollection)
// 		if err != nil {
// 			log.Fatalln(err)
// 		}
//
// 		cmd, err := leaf.NewCommand("npm run build")
// 		if err != nil {
// 			log.Fatalln(err)
// 		}
//
// 		log.Printf("Watching: %s\n", root)
//
// 		for change := range watcher.Watch(ctx) {
// 			if change.Err != nil {
// 				log.Printf("ERROR: %v", change.Err)
// 				continue
// 			}
// 			// If no error run the command
// 			log.Printf("Running: %s\n", cmd.String())
// 			cmd.Execute(ctx)
// 		}
// 	}
//
package leaf

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"runtime/debug"
	"time"

	"github.com/sirupsen/logrus"
)

var (
	// DefaultExcludePathsKeyword is used to include all
	// default excludes.
	DefaultExcludePathsKeyword = "DEFAULTS"

	// CWD is the current working directory or ".".
	CWD string

	// DefaultConfPath is the default path for app config.
	DefaultConfPath string

	// DefaultExcludePaths are the paths that should be
	// generally excluded while watching a project.
	DefaultExcludePaths = []string{
		".git/",
		"node_modules/",
		"vendor/",
		"venv/",
	}
	// ImportPath is the import path for leaf package.
	ImportPath = "github.com/vrongmeal/leaf"
)

func init() {
	var err error
	CWD, err = os.Getwd()
	if err != nil {
		logrus.Fatalln(err)
	}

	DefaultConfPath = filepath.Join(CWD, ".leaf.yml")
}

// Config represents the conf file for the runner.
type Config struct {
	// Root directory to watch.
	Root string `mapstructure:"root"`

	// Exclude these directories from watch.
	Exclude []string `mapstructure:"exclude"`

	// Filters to apply to the watch.
	Filters []string `mapstructure:"filters"`

	// Exec these commads after changes detected.
	Exec []string `mapstructure:"exec"`

	// ExitOnErr breaks the chain of command if any command returnns an error.
	ExitOnErr bool `mapstructure:"exit_on_err"`

	// Delay after which commands should be executed.
	Delay time.Duration `mapstructure:"delay"`
}

// NewCmdContext returns a context which cancels on an OS
// interrupt, i.e., cancels when process is killed.
func NewCmdContext(onInterrupt func(os.Signal)) context.Context {
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

// GoModuleInfo returns the go module information which
// includes the build info (version etc.).
func GoModuleInfo() (*debug.Module, error) {
	buildInfo, ok := debug.ReadBuildInfo()
	if !ok {
		return nil, fmt.Errorf("unable to fetch build info")
	}

	for _, dep := range buildInfo.Deps {
		if dep.Path == ImportPath {
			return dep, nil
		}
	}

	return &buildInfo.Main, nil
}

// *** some helper functions ***

// isDir checks if the given path is a directory or not.
// Returns an error when path is invalid.
func isDir(root string) (bool, error) {
	info, err := os.Stat(root)
	if err != nil {
		if os.IsNotExist(err) {
			return false, fmt.Errorf(
				"filepath does not exist")
		}

		return false, err
	}

	if info.IsDir() {
		return true, nil
	}

	return false, nil
}
