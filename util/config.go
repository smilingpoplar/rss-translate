package util

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Output OutputConfig          `yaml:"output"`
	ToLang string                `yaml:"to-lang"`
	Proxy  string                `yaml:"proxy"`
	Feeds  map[string]FeedConfig `yaml:"feeds"`
	Glossary map[string]string    `yaml:"glossary"`
}

type OutputConfig struct {
	Dir string `yaml:"dir"`
	URL string `yaml:"url"`
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
