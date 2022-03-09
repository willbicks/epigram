package config

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

// UnmarshallYAML chooses the apporopriate constant value for an entry of type Repository
func (repo *Repository) UnmarshalYAML(n *yaml.Node) error {
	var s string

	if err := n.Decode(&s); err != nil {
		return fmt.Errorf("decoding value for repository: %w", err)
	}

	switch strings.ToLower(s) {
	case "inmemory":
		*repo = Inmemory
	case "sqlite":
		*repo = SQLite
	default:
		return fmt.Errorf("unexpected repository value '%v'", s)
	}

	return nil
}

// ParseYAML accepts an array of bytes, and parses it into an Applicaiton configuration
func ParseYAML(in []byte) (Application, error) {
	var cfg Application

	if err := yaml.Unmarshal(in, &cfg); err != nil {
		return Application{}, err
	}

	return cfg, nil
}
