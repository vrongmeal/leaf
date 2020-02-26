package utils

import (
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

var (
	// DefaultExcludePathsKeyword is used to include all default excludes.
	DefaultExcludePathsKeyword = "DEFAULTS"

	// CWD is the current working directory or ".".
	CWD string

	// DefaultConfPath is the default path for app config.
	DefaultConfPath string

	// DefaultExcludePaths are the paths that should be generally excluded
	// while watching a project.
	DefaultExcludePaths = []string{
		".git/",
		"node_modules/",
		"vendor/",
		"venv/",
	}
)

func init() {
	var err error
	CWD, err = os.Getwd()
	if err != nil {
		logrus.Fatalln(err)
	}

	DefaultConfPath = filepath.Join(CWD, ".leaf.yml")
}
