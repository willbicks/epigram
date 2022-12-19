package config

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

// UnmarshalYAML chooses the appropriate constant value for an entry of type Repository
func (repo *Repository) UnmarshalYAML(n *yaml.Node) error {
	var str string

	if err := n.Decode(&str); err != nil {
		return fmt.Errorf("decoding value for repository: %w", err)
	}

	*repo = repoFromString(str)

	return nil
}

// parseYAML accepts an array of bytes, and parses it into an Application configuration
func parseYAML(in []byte) (Application, error) {
	var cfg Application

	if err := yaml.Unmarshal(in, &cfg); err != nil {
		return Application{}, err
	}

	return cfg, nil
}
