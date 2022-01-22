package configutil

import (
	"fmt"
	"os"

	"github.com/caarlos0/env/v6"
	"gopkg.in/yaml.v2"
)

func LoadConfigFromFile(configPath string, cfg interface{}) error {
	content, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to read config %s: %w", configPath, err)
	}

	err = yaml.Unmarshal(content, cfg)
	if err != nil {
		return fmt.Errorf("failed to read yaml: %w", err)
	}

	return nil
}

func LoadConfigFromEnv(cfg interface{}) error {
	if err := env.Parse(cfg); err != nil {
		return fmt.Errorf("failed to init app config: %w", err)
	}

	return nil
}
