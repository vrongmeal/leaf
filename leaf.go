// Package leaf provides with utilities to create the leaf
// CLI tool. It includes watcher, filters and commander which
// watch files for changes, filter out required results and
// execute external commands respectively.
package leaf

import (
	"fmt"
	"os"
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

	// Delay after which commands should be executed.
	Delay time.Duration `mapstructure:"delay"`
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
