package watcher

import (
	"fmt"
	"path/filepath"
	"strings"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/sirupsen/logrus"

	"github.com/vrongmeal/leaf/pkg/utils"
)

// Watcher is the data type that handles watching of paths.
type Watcher struct {
	root    string
	paths   []string
	exclude []string

	*filterCollection

	notifier *fsnotify.Watcher

	done chan bool
	wg   *sync.WaitGroup

	File chan string
	Err  chan error
}

// WatchOpts are the options used to configure the watcher.
type WatchOpts struct {
	Root    string
	Exclude []string
	Filters []string
}

// NewWatcher returns a watcher from the given options.
func NewWatcher(opts *WatchOpts) (*Watcher, error) {
	w := &Watcher{}

	isdir, err := utils.IsDir(opts.Root)
	if err != nil {
		return nil, err
	}

	if !isdir {
		return nil, fmt.Errorf("path '%s' is not a directory", opts.Root)
	}

	w.root = filepath.Clean(opts.Root)
	w.paths, err = utils.GetAllDirs(w.root)
	if err != nil {
		return nil, err
	}

	w.exclude = []string{}
	for _, path := range opts.Exclude {
		var absPath string

		if _, err = utils.IsDir(path); err != nil {
			continue
		}

		absPath, err = filepath.Abs(path)
		if err != nil {
			continue
		}

		w.exclude = append(w.exclude, absPath)
	}

	w.filterCollection = newFilterCollection(opts.Filters, func(err error) { logrus.Errorln(err) })

	w.notifier, err = fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	w.done = make(chan bool)
	w.wg = &sync.WaitGroup{}

	w.File = make(chan string)
	w.Err = make(chan error)

	return w, nil
}

// Watch executes the watching of files.
func (w *Watcher) Watch() error {
	defer w.notifier.Close() // nolint:errcheck

	w.wg.Add(1)
	go func() {
		defer w.wg.Done()
		for {
			select {
			case event := <-w.notifier.Events:
				if event.Op == fsnotify.Write {
					file := event.Name
					handle := false
					if len(w.includes) == 0 || w.hasInclude(file) {
						handle = true
					}
					if w.hasExclude(file) {
						handle = false
					}
					if handle {
						w.File <- file
					}
				}
			case err := <-w.notifier.Errors:
				if err != nil {
					w.Err <- err
				}
			case <-w.done:
				return
			}
		}
	}()

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
			w.done <- true
			w.wg.Wait()
			return err
		}
	}

	w.wg.Wait()

	return nil
}

// Close terminates the watcher.
func (w *Watcher) Close() {
	w.done <- true
	w.wg.Wait()
}
