# leaf

> General purpose hot-reloader for all projects.

## Contents

1. [Installation](#installation)
  1. [Using `go get`](#using-go-get)
  1. [Manual](#manual)
1. [Usage](#usage)
  1. [Command line help](#command-line-help)
  1. [Configuration file](#configuration-file)

## Installation

### Using `go get`

The following command will download and build Leaf in your `$GOPATH/bin`.

```console
> go get -u github.com/vrongmeal/leaf
```

### Manual

1. Clone the repository and `cd` into it.
1. Run `make build` to build the leaf as `build/leaf`.
1. Move the binary somewhere in your `$PATH`.

## Usage

### Command line help

```console
> leaf --help
Given a set of commands, leaf watches the filtered paths in the project directory for any changes and runs the commands in
order so you don't have to yourself

Usage:
  leaf [flags]
  leaf [command]

Available Commands:
  help        Help about any command
  version     Leaf version

Flags:
  -c, --config string   Config path for leaf configuration file (default "<CWD>/.leaf.yml")
  -h, --help            help for leaf

Use "leaf [command] --help" for more information about a command.
```

### Configuration file

This project doesn't really require a hot-reload but a sample hot-reload for this project would look like [this](_examples/sample.leaf.yml):

```yaml
# Leaf configuration file.

# Root directory to watch.
# Defaults to current working directory.
root: "."

# Exclude directories while watching.
# If certain directories are not excluded, it might reach a limitation where watcher doesn't start.
exclude:
  - ".git/"
  - "vendor/"
  - "build/"

# Filters to apply on the watch.
# Filters starting with '+' are includent and then with '-' are excluded.
# This is not like exclude, these are still being watched yet can be excluded from the execution.
# These can include any filepath regex supported by "filepath".Match method or even a directory.
filters:
  - "- .git*"
  - "- .go*"
  - "- .golangci.yml"
  - "- go.*"
  - "- Makefile"
  - "- LICENSE"
  - "- README.md"

# Commands to be executed.
# These are run in the provided order.
exec:
  - ["make", "build"]

# Delay after which commands are executed.
delay: '1s'
```

By default the config path is taken as `<current working directory>/.leaf.yml` which you can change using the `--config` or `-c` flag. You can also use a JSON or TOML file. Just remember to use "delay" in nanoseconds (seconds * 10^9).

---

Made with **khoon**, **paseena** and **love** `:-)` by

Vaibhav ([vrongmeal](https://vrongmeal.github.io))
