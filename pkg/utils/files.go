package utils

import (
	"errors"
	"os"
	"path/filepath"
)

var (
	// ErrorInvalidPath is returned when filepath is invalid.
	ErrorInvalidPath = errors.New("filepath does not exist")
)

// GetAllDirs gets all the directories (including the root) inside the given root directory.
func GetAllDirs(root string) ([]string, error) {
	paths := []string{}

	if err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			absPath, err := filepath.Abs(path)
			if err != nil {
				return err
			}

			paths = append(paths, absPath)
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return paths, nil
}

// IsDir checks if the given path is a directory or not. Returns an error when path is invalid.
func IsDir(root string) (bool, error) {
	info, err := os.Stat(root)
	if err != nil {
		if os.IsNotExist(err) {
			return false, ErrorInvalidPath
		}

		return false, err
	}

	if info.IsDir() {
		return true, nil
	}

	return false, nil
}
