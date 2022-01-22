package configutil

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

func LoadConfig(configPath string, cfg interface{}) error {
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
