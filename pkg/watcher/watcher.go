package watcher

import (
	"fmt"
	"path/filepath"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/sirupsen/logrus"

	"github.com/vrongmeal/leaf/pkg/utils"
)

// Watcher is the data type that handles watching of paths.
type Watcher struct {
	root  string
	paths []string

	*filterCollection

	notifier *fsnotify.Watcher

	fileFn FileFunc
	errFn  ErrorFunc
	stopFn StopFunc

	done chan bool
	wg   *sync.WaitGroup
}

// WatchOpts are the options used to configure the watcher.
type WatchOpts struct {
	Root    string
	Filters []string

	FileFn  FileFunc
	ErrorFn ErrorFunc
	StopFn  StopFunc
}

// FileFunc is the function that runs when a file is received by the channel.
type FileFunc func(file string)

func nilFileFunc(string) {}

// ErrorFunc is the function that runs when an error occurs while watching.
type ErrorFunc func(err error)

func nilErrorFunc(err error) { logrus.Errorln(err) }

// StopFunc is the function that runs when the watcher is closed.
type StopFunc func()

func nilStopFunc() {}

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

	w.filterCollection = newFilterCollection(opts.Filters, nilErrorFunc)

	if opts.FileFn == nil {
		opts.FileFn = nilFileFunc
	}

	if opts.ErrorFn == nil {
		opts.ErrorFn = nilErrorFunc
	}

	if opts.StopFn == nil {
		opts.StopFn = nilStopFunc
	}

	w.notifier, err = fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	w.fileFn = opts.FileFn
	w.errFn = opts.ErrorFn
	w.stopFn = opts.StopFn

	w.done = make(chan bool)
	w.wg = &sync.WaitGroup{}

	return w, nil
}

// Watch executes the watching of files.
func (w *Watcher) Watch() error {
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
						w.handleFile(event.Name)
					}
				}
			case err := <-w.notifier.Errors:
				if err != nil {
					w.handleErr(err)
				}
			case <-w.done:
				w.handleClose()
				return
			}
		}
	}()

	for _, f := range w.paths {
		if err := w.notifier.Add(f); err != nil {
			w.done <- true
			w.wg.Wait()
			return err
		}
	}

	return nil
}

// Close terminates the watcher.
func (w *Watcher) Close() error {
	w.done <- true
	w.wg.Wait()
	return w.notifier.Close()
}

// Wait waits for the watch to close.
func (w *Watcher) Wait() {
	w.wg.Wait()
}

func (w *Watcher) handleFile(file string) {
	fileFn := w.fileFn
	fileFn(file)
}

func (w *Watcher) handleErr(err error) {
	errFn := w.errFn
	errFn(err)
}

func (w *Watcher) handleClose() {
	stopFn := w.stopFn
	stopFn()
}
