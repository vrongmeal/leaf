# leaf

> General purpose hot-reloader for all projects.

## Contents

1. [Installation](#installation)
    1. [Homebrew](#homebrew)
    1. [Using `go get`](#using-go-get)
    1. [Manual](#manual)
1. [Usage](#usage)
    1. [Command line help](#command-line-help)
    1. [Configuration file](#configuration-file)
    1. [Configuring using cmd](#configuring-using-cmd)

## Installation

### Homebrew

You can use my homebrew tap to install Leaf.

```console
> brew tap vrongmeal/tap
> brew install vrongmeal/tap/leaf
```

### Using `go get`

The following command will download and build Leaf in your `$GOPATH/bin`.

```console
> go get -u github.com/vrongmeal/leaf
```

**Note:** This does not build the `version` command. To build that use [homebrew](#homebrew) or [manual](#manual) installation.

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
  -c, --config string      Config path for leaf configuration file (default "<CWD>/.leaf.yml")
  -d, --delay duration     Delay after which commands are run on file change (default 500ms)
  -e, --exclude DEFAULTS   Paths to exclude. You can append default paths by adding DEFAULTS in your list (default [.git/,node_modules/,vendor/,venv/])
  -x, --exec strings       Exec commands on file change
  -f, --filters strings    Filters to apply to watch
  -h, --help               help for leaf
  -r, --root string        Root directory to watch (default "<CWD>")

Use "leaf [command] --help" for more information about a command.
```

### Configuration file

This project doesn't really require a hot-reload but a sample hot-reload for this project would look like [this](_examples/sample.leaf.yml):

```yaml
# Leaf configuration file.

# Root directory to watch.
# Defaults to current working directory.
root: .

# Exclude directories while watching.
# If certain directories are not excluded, it might reach a limitation where watcher doesn't start.
exclude:
  - DEFAULTS # This includes the default ignored directories
  - build/

# Filters to apply on the watch.
# Filters starting with '+' are includent and then with '-' are excluded.
# This is not like exclude, these are still being watched yet can be excluded from the execution.
# These can include any filepath regex supported by "filepath".Match method or even a directory.
filters:
  # The following can be simplified by also doing an include filter:
  # ['+cmd/', '+pkg/', '+scripts/']
  # This example is just to show that expressions are supported.
  - -.git*
  - -.go*
  - -.golangci.yml
  - -go.*
  - -Makefile
  - -LICENSE
  - -README.md

# Commands to be executed.
# These are run in the provided order.
exec:
  - make build

# Delay after which commands are executed.
delay: 1s
```

By default the config path is taken as `<current working directory>/.leaf.yml` which you can change using the `--config` or `-c` flag. You can also use a JSON or TOML file.

## Configuring using cmd

You can also configure using command line. The above config can be run as follows:

```console
> leaf -d=1s -e=DEFAULTS -e='build/' -f='+cmd/' -f='+pkg/' -f='+scripts/' -x='make build'
[0000]  WARN Continuing without config...
[0000]  INFO Starting to watch: /Users/vrongmeal/Projects/vrongmeal/leaf
[0000]  INFO Excluded paths:
[0000]  INFO  .git/
[0000]  INFO  node_modules/
[0000]  INFO  vendor/
[0000]  INFO  venv/
[0000]  INFO  build/
...
```

---

Made with **khoon**, **paseena** and **love** `:-)` by

Vaibhav ([vrongmeal](https://vrongmeal.github.io))
