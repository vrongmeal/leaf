package leaf

import (
	"fmt"
	"path/filepath"
	"strings"
)

// Filter can be used to Filter out watch results.
type Filter struct {
	include bool
	pattern string
}

// NewFilter creates a filter from the pattern string. The
// pattern either starts with '+' or '-' to include or
// exclude the directory from results.
func NewFilter(pattern string) (Filter, error) {
	f := Filter{}
	var err error

	cleanedPattern := strings.Trim(pattern, " ")
	if len(cleanedPattern) < 2 {
		return f, fmt.Errorf(
			"effective pattern '%s' invalid", cleanedPattern)
	}

	toInclude := cleanedPattern[0]
	if toInclude == '+' {
		f.include = true
	} else if toInclude == '-' {
		f.include = false
	} else {
		return f, fmt.Errorf(
			"should have first character as '+' or '-'")
	}

	onlyPath := strings.Trim(cleanedPattern[1:], " ")
	f.pattern, err = filepath.Abs(onlyPath)
	if err != nil {
		return f, fmt.Errorf(
			"error making path absolute: %v", err)
	}

	return f, nil
}

// A FilterCollection contains a bunch of includes and excludes.
type FilterCollection struct {
	Includes []Filter
	Excludes []Filter
}

// NewFilterCollection creates a filter collection from a bunch
// of filter patterns.
func NewFilterCollection(patterns []string) (*FilterCollection, error) {
	collection := &FilterCollection{
		Includes: []Filter{},
		Excludes: []Filter{},
	}

	if len(patterns) == 0 {
		return collection, nil
	}

	for _, pattern := range patterns {
		f, err := NewFilter(pattern)
		if err != nil {
			return nil, fmt.Errorf(
				"error in '%s' pattern: %v", pattern, err)
		}

		if f.include {
			collection.Includes = append(collection.Includes, f)
		} else {
			collection.Excludes = append(collection.Excludes, f)
		}
	}

	return collection, nil
}

// match matches the pattern with the path and returns true
// if the path either starts with (in absolute terms) or
// matches like the path regex.
func match(pattern, path string) bool {
	matched, err := filepath.Match(pattern, path)
	if err != nil {
		return false
	}

	if matched {
		return true
	}

	isDir, err := isDir(pattern)
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

// HasInclude tells if the collection matches the path with
// one of its includes.
func (c *FilterCollection) HasInclude(path string) bool {
	cleanedPath := filepath.Clean(path)

	for _, f := range c.Includes {
		if match(f.pattern, cleanedPath) {
			return true
		}
	}

	return false
}

// HasExclude tells if the collection matches the path with
// one of its excludes.
func (c *FilterCollection) HasExclude(path string) bool {
	cleanedPath := filepath.Clean(path)

	for _, f := range c.Excludes {
		if match(f.pattern, cleanedPath) {
			return true
		}
	}

	return false
}
