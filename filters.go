package leaf

import (
	"fmt"
	"path/filepath"
	"strings"
)

// Filter can be used to Filter out watch results.
type Filter struct {
	Include bool // whether to include pattern
	Pattern string
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
		f.Include = true
	} else if toInclude == '-' {
		f.Include = false
	} else {
		return f, fmt.Errorf(
			"should have first character as '+' or '-'")
	}

	onlyPath := strings.Trim(cleanedPattern[1:], " ")
	f.Pattern, err = filepath.Abs(onlyPath)
	if err != nil {
		return f, fmt.Errorf(
			"error making path absolute: %v", err)
	}

	return f, nil
}

// A FilterCollection contains a bunch of includes and excludes.
type FilterCollection struct {
	Includes []string
	Excludes []string

	match  FilterMatchFunc
	handle FilterHandleFunc
}

// NewFilterCollection creates a filter collection from a bunch
// of filter patterns.
func NewFilterCollection(filters []Filter, mf FilterMatchFunc, hf FilterHandleFunc) *FilterCollection {
	collection := &FilterCollection{
		Includes: []string{},
		Excludes: []string{},
		match:    mf,
		handle:   hf,
	}

	if len(filters) == 0 {
		return collection
	}

	for _, f := range filters {
		if f.Include {
			collection.Includes = append(collection.Includes, f.Pattern)
		} else {
			collection.Excludes = append(collection.Excludes, f.Pattern)
		}
	}

	return collection
}

// NewFCFromPatterns creates a filter collection from a list of
// string format filters, like `+ /path/to/some/dir`.
func NewFCFromPatterns(patterns []string, mf FilterMatchFunc, hf FilterHandleFunc) (*FilterCollection, error) {
	filters := []Filter{}

	for _, p := range patterns {
		f, err := NewFilter(p)
		if err != nil {
			return nil, err
		}

		filters = append(filters, f)
	}

	return NewFilterCollection(filters, mf, hf), nil
}

// FilterMatchFunc compares the pattern with the path of
// the file changed and returns true if the path resembles
// the given pattern.
type FilterMatchFunc func(pattern, path string) bool

// StandardFilterMatcher matches the pattern with the path
// and returns true if the path either starts with
// (in absolute terms) or matches like the path regex.
func StandardFilterMatcher(pattern, path string) bool {
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
func (fc *FilterCollection) HasInclude(path string) bool {
	cleanedPath := filepath.Clean(path)

	for _, pattern := range fc.Includes {
		if fc.match(pattern, cleanedPath) {
			return true
		}
	}

	return false
}

// HasExclude tells if the collection matches the path with
// one of its excludes.
func (fc *FilterCollection) HasExclude(path string) bool {
	cleanedPath := filepath.Clean(path)

	for _, pattern := range fc.Excludes {
		if fc.match(pattern, cleanedPath) {
			return true
		}
	}

	return false
}

// ShouldHandlePath returns the result of the path handler
// for the filter collection.
func (fc *FilterCollection) ShouldHandlePath(path string) bool {
	handlerFunc := fc.handle
	return handlerFunc(fc, path)
}

// FilterHandleFunc is a function that checks if for the filter
// collection, should the path be handled or not, i.e., should
// the notifier tick for change in path or not.
type FilterHandleFunc func(fc *FilterCollection, path string) bool

// StandardFilterHandler returns true if the path should be included
// and returns false if path should not be included in result.
func StandardFilterHandler(fc *FilterCollection, path string) bool {
	handle := false

	// If there are no includes, path should be handled unless
	// it is in the excludes.
	if len(fc.Includes) == 0 || fc.HasInclude(path) {
		handle = true
	}

	if fc.HasExclude(path) {
		handle = false
	}

	return handle
}
