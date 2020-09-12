# leaf

> General purpose hot-reloader for all projects.

![Continuous Integration](https://github.com/vrongmeal/leaf/workflows/Continuous%20Integration/badge.svg)

Command leaf watches for changes in the working directory
and runs the specified set of commands whenever a file
updates. A set of filters can be applied to the watch and
directories can be excluded.

## Contents

1. [Installation](#installation)
    1. [Using `go get`](#using-go-get)
    1. [Manual](#manual)
1. [Usage](#usage)
    1. [Command line help](#command-line-help)
    1. [Configuration file](#configuration-file)
1. [Custom hot reloader](#custom-hot-reloader)

## Installation

### Using `go get`

The following command will download and build Leaf in your
`$GOPATH/bin`.

```
❯ go get -u github.com/vrongmeal/leaf/cmd/leaf
```

### Manual

1. Clone the repository and `cd` into it.
1. Run `make build` to build the leaf as `build/leaf`.
1. Move the binary somewhere in your `$PATH`.

## Usage

```
❯ leaf -x 'make build' -x 'make run'
```

The above command runs `make build` and `make run` commands
(in order).

### Command line help

The CLI can be used as described by the help message:

```
❯ leaf help

Command leaf watches for changes in the working directory and
runs the specified set of commands whenever a file updates.
A set of filters can be applied to the watch and directories
can be excluded.

Usage:
  leaf [flags]
  leaf [command]

Available Commands:
  help        Help about any command
  version     prints leaf version

Flags:
  -c, --config string     config path for the configuration file (default "<CWD>/.leaf.yml")
      --debug             run in development (debug) environment
  -d, --delay duration    delay after which commands are run on file change (default 500ms)
  -e, --exclude strings   paths to exclude from watching (default [.git/,node_modules/,vendor/,venv/])
  -x, --exec strings      exec commands on file change
  -z, --exit-on-err       exit chain of commands on error
  -f, --filters strings   filters to apply to watch
  -h, --help              help for leaf
  -o, --once              run once and exit (no reload)
  -r, --root string       root directory to watch (default "<CWD>")

Use "leaf [command] --help" for more information about a command.
```

### Configuration file

In order to configure using a configuration file, create a
YAML or TOML or even a JSON file with the following structure
and pass it using the `-c` or `--config` flag. By default
a file named `.leaf.yml` in your working directory is taken
if no configuration file is found.

```yaml
# Leaf configuration file.

# Root directory to watch.
# Defaults to current working directory.
root: .

# Exclude directories while watching.
# If certain directories are not excluded, it might reach a
# limitation where watcher doesn't start.
exclude:
  - DEFAULTS # This includes the default ignored directories
  - build/
  - scripts/

# Filters to apply on the watch.
# Filters starting with '+' are includent and then with '-'
# are excluded. This is not like exclude, these are still
# being watched yet can be excluded from the execution.
# These can include any regex supported by filepath.Match
# method or even a directory.
filters:
  - '+ go.mod'
  - '+ go.sum'
  - '+ *.go'
  - '+ cmd/'

# Commands to be executed. These are run in the provided order.
exec:
  - make format
  - make build

# Stop the command chain when an error occurs
exit_on_err: true

# Delay after which commands are executed.
delay: 1s
```

The above config file is suitable to use with the current
project itself. It can also be translated into a command
as such:

```
❯ leaf -z -x 'make format' -x 'make build' -d '1s' \
  -e 'DEFAULTS' -e 'build' -e 'scripts' \
  -f '+ go.*' -f '+ *.go' -f '+ cmd/'
```

## Custom hot reloader

The package [github.com/vrongmeal/leaf](https://pkg.go.dev/github.com/vrongmeal/leaf)
comes with utilities that can aid in creating a hot-reloader
with a simple go program.

Let's look at an example where the watcher watches the `src/`
directory for changes and for any changes builds the project.

```go
package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/vrongmeal/leaf"
)

func main() {
	// Use a context that cancels when program is interrupted.
	ctx := leaf.NewCmdContext(func(os.Signal) {
		log.Println("Shutting down.")
	})

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalln(err)
	}

	// Root is <cwd>/src
	root := filepath.Join(cwd, "src")

	// Exclude "src/index.html" from results.
	filters := []leaf.Filter{
		{Include: false, Pattern: "src/index.html"},
	}

	filterCollection := leaf.NewFilterCollection(
		filters,
		// Matches directory or filepath.Match expressions
		leaf.StandardFilterMatcher,
		// Definitely excludes and shows only includes (if any)
		leaf.StandardFilterHandler)

	watcher, err := leaf.NewWatcher(
		root,
		// Standard paths to exclude, like vendor, .git,
		// node_modules, venv etc.
		leaf.DefaultExcludePaths,
		filterCollection)
	if err != nil {
		log.Fatalln(err)
	}

	cmd, err := leaf.NewCommand("npm run build")
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("Watching: %s\n", root)

	for change := range watcher.Watch(ctx) {
		if change.Err != nil {
			log.Printf("ERROR: %v", change.Err)
			continue
		}
		// If no error run the command
		fmt.Printf("Running: %s\n", cmd.String())
		cmd.Execute(ctx)
	}
}
```

---

Made with **khoon**, **paseena** and **love** `:-)` by

Vaibhav ([vrongmeal](https://vrongmeal.github.io))
