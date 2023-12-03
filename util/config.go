package util

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Output OutputConfig          `yaml:"output"`
	Feeds  map[string]FeedConfig `yaml:"feeds"`
}

type OutputConfig struct {
	Dir string `yaml:"dir"`
}

type FeedConfig struct {
	URL string `yaml:"url"`
	Max int    `yaml:"max"`
}

func GetConfig(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("error opening config file %s: %w", path, err)
	}
	defer f.Close()

	config := Config{}
	err = yaml.NewDecoder(f).Decode(&config)
	if err != nil {
		return nil, fmt.Errorf("error decoding config: %w", err)
	}

	return &config, nil
}
