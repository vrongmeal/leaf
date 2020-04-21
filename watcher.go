package leaf

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
)

// WatchResult has the file changed or the error that occurred
// during watching.
type WatchResult struct {
	File string
	Err  error
}

// Watcher watches a directory for changes and updates the
// stream when a file change (valid by filters) is updated.
type Watcher struct {
	root    string
	paths   []string
	exclude []string

	fc       *FilterCollection
	notifier *fsnotify.Watcher

	res chan WatchResult
}

// NewWatcher returns a watcher from the given options.
func NewWatcher(root string, exclude, filters []string) (*Watcher, error) {
	w := &Watcher{}

	isdir, err := isDir(root)
	if err != nil {
		return nil, err
	}

	if !isdir {
		return nil, fmt.Errorf("path '%s' is not a directory", root)
	}

	w.root = filepath.Clean(root)
	w.paths, err = getAllDirs(w.root)
	if err != nil {
		return nil, err
	}

	w.exclude = []string{}
	for _, path := range exclude {
		var absPath string

		if _, err = isDir(path); err != nil {
			continue
		}

		absPath, err = filepath.Abs(path)
		if err != nil {
			continue
		}

		w.exclude = append(w.exclude, absPath)
	}

	w.fc, err = NewFilterCollection(filters)
	if err != nil {
		return nil, err
	}

	w.notifier, err = fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	for _, f := range w.paths {
		exclude := false
		for _, e := range w.exclude {
			if strings.HasPrefix(f, e) {
				exclude = true
				break
			}
		}

		if exclude {
			continue
		}

		if err := w.notifier.Add(f); err != nil {
			return nil, err
		}
	}

	w.res = make(chan WatchResult)

	return w, nil
}

// Watch executes the watching of files. Exits on cancellation
// of the context.
func (w *Watcher) Watch(ctx context.Context) <-chan WatchResult {
	go w.startWatcher(ctx)
	return w.res
}

// startWatcher starts the fs.Notifier and watches for changes
// in files in the root directory.
func (w *Watcher) startWatcher(ctx context.Context) {
	defer w.notifier.Close() // nolint:errcheck
	for {
		select {
		case event := <-w.notifier.Events:
			if event.Op == fsnotify.Write {
				file := event.Name
				handle := false
				if len(w.fc.Includes) == 0 || w.fc.HasInclude(file) {
					handle = true
				}

				if w.fc.HasExclude(file) {
					handle = false
				}

				if handle {
					w.res <- WatchResult{File: file}
				}
			}

		case err := <-w.notifier.Errors:
			if err != nil {
				w.res <- WatchResult{Err: err}
			}

		case <-ctx.Done():
			close(w.res)
			return
		}
	}
}

// getAllDirs gets all the directories (including the root)
// inside the given root directory.
func getAllDirs(root string) ([]string, error) {
	paths := []string{}

	walkFn := func(path string, info os.FileInfo, err error) error {
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
	}

	err := filepath.Walk(root, walkFn)
	if err != nil {
		return nil, err
	}

	return paths, nil
}
