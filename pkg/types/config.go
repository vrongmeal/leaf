package types

import (
	"time"
)

// Config represents the conf file for the runner.
type Config struct {
	Root    string        `json:"root" yaml:"root" toml:"root"`
	Exclude []string      `json:"exclude" yaml:"exclude" toml:"exclude"`
	Filters []string      `json:"filters" yaml:"filters" toml:"filters"`
	Exec    []string      `json:"exec" yaml:"exec" toml:"exec"`
	Delay   time.Duration `json:"delay" yaml:"delay" toml:"delay"`
}
