package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

func ReadFromFile(filePath string) (*Config, error) {
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path: %w", err)
	}

	data, err := os.ReadFile(absPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("%w: %s", ErrFileNotFound, absPath)
		}
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidYAML, err)
	}

	for i, operator := range config.Operators {
		if err := operator.Validate(); err != nil {
			return nil, fmt.Errorf("invalid operator '%s' at index %d: %w", operator.Name, i, err)
		}
	}

	return &config, nil
}
