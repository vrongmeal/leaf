package utils

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"time"

	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v2"
)

// Config represents the conf file for the runner.
type Config struct {
	Root    string        `json:"root" yaml:"root" toml:"root"`
	Exclude []string      `json:"exclude" yaml:"exclude" toml:"exclude"`
	Filters []string      `json:"filters" yaml:"filters" toml:"filters"`
	Exec    [][]string    `json:"exec" yaml:"exec" toml:"exec"`
	Delay   time.Duration `json:"delay" yaml:"delay" toml:"delay"`
}

// GetConfig returns config from the filepath.
func GetConfig(path string) (Config, error) {
	config := Config{}

	isdir, err := IsDir(path)
	if err != nil || isdir {
		return config, err
	}

	content, err := ioutil.ReadFile(filepath.Clean(path))
	if err != nil {
		return config, err
	}

	switch filepath.Ext(path) {
	case ".json":
		if err := json.Unmarshal(content, &config); err != nil {
			return config, err
		}
	case ".toml":
		if _, err := toml.Decode(string(content), &config); err != nil {
			return config, err
		}
	default:
		if err := yaml.Unmarshal(content, &config); err != nil {
			return config, err
		}
	}

	return config, nil
}
