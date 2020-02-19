package watcher

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/vrongmeal/leaf/pkg/utils"
)

type filter struct {
	include bool
	pattern string
}

func newFilter(pattern string) (filter, error) {
	f := filter{}
	var err error

	cleanedPattern := strings.Trim(pattern, " ")
	if len(cleanedPattern) < 2 {
		return f, fmt.Errorf("invalid filter pattern '%s'", cleanedPattern)
	}

	toInclude := cleanedPattern[0]
	if toInclude == '+' {
		f.include = true
	} else if toInclude == '-' {
		f.include = false
	} else {
		return f, fmt.Errorf("invalid filter pattern: should have the first character as '+' or '-'")
	}

	onlyPath := strings.Trim(cleanedPattern[1:], " ")
	f.pattern, err = filepath.Abs(onlyPath)
	if err != nil {
		return f, fmt.Errorf("pattern could not be made absolute")
	}

	return f, nil
}

type filterCollection struct {
	includes []filter
	excludes []filter
}

func newFilterCollection(patterns []string, onErrorFn func(error)) *filterCollection {
	collection := &filterCollection{
		includes: []filter{},
		excludes: []filter{},
	}

	if len(patterns) == 0 {
		return collection
	}

	for _, pattern := range patterns {
		f, err := newFilter(pattern)
		if err != nil {
			onErrorFn(err)
			continue
		}

		if f.include {
			collection.includes = append(collection.includes, f)
		} else {
			collection.excludes = append(collection.excludes, f)
		}
	}

	return collection
}

func match(pattern, path string) bool {
	matched, err := filepath.Match(pattern, path)
	if err != nil {
		return false
	}

	if matched {
		return true
	}

	isDir, err := utils.IsDir(pattern)
	if err != nil || !isDir {
		return false
	}

	if len(path) < len(pattern) {
		return false
	}

	if path[:len(pattern)] == pattern {
		return true
	}

	return false
}

func (c *filterCollection) hasInclude(path string) bool {
	cleanedPath := filepath.Clean(path)

	for _, f := range c.includes {
		if match(f.pattern, cleanedPath) {
			return true
		}
	}

	return false
}

func (c *filterCollection) hasExclude(path string) bool {
	cleanedPath := filepath.Clean(path)

	for _, f := range c.excludes {
		if match(f.pattern, cleanedPath) {
			return true
		}
	}

	return false
}
