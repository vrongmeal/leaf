package utils

import (
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

var (
	// CWD is the current working directory or "."
	CWD string

	// DefaultConfPath is the default path for app config.
	DefaultConfPath string
)

func init() {
	var err error
	CWD, err = os.Getwd()
	if err != nil {
		logrus.Fatalln(err)
	}

	DefaultConfPath = filepath.Join(CWD, ".leaf.yml")
}
