package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Target struct {
	Name string `yaml:"name"`
	URL  string `yaml:"url"`
}

type Telegram struct {
	Enabled  bool   `yaml:"enabled"`
	BotToken string `yaml:"bot_token"`
	ChatID   string `yaml:"chat_id"`
}

type Config struct {
	CheckIntervalSeconds  int      `yaml:"check_interval_seconds"`
	RequestTimeoutSeconds int      `yaml:"request_timeout_seconds"`
	Telegram              Telegram `yaml:"telegram"`
	Targets               []Target `yaml:"targets"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}