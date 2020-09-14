// Package matcher matches with file paths using the gitignore syntax.
package matcher

import (
	"path/filepath"

	ignore "github.com/sabhiram/go-gitignore"
)

// Matcher matches the gitignore patterns with the file paths.
type Matcher struct{ g *ignore.GitIgnore }

// NewMatcher creates a new matcher for checking if filepaths match given
// patterns according to the gitignore syntax.
func NewMatcher(patterns []string) (Matcher, error) {
	gign, err := ignore.CompileIgnoreLines(patterns...)
	if err != nil {
		return Matcher{}, err
	}

	return Matcher{g: gign}, nil
}

// Match checks if the file path matches the patterns.
func (m Matcher) Match(path string, isdir bool) bool {
	path = filepath.Clean(path)
	if path == "" {
		return false
	}

	if isdir && path[len(path)-1] != '/' {
		path = path + "/"
	}

	return m.g.MatchesPath(path)
}
