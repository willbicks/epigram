// Package config is responsible for processing global, application configuration specified at the time of execution.
// Supports both reading from .yml files, and environment variables, the latter taking priority over the former.
//
// For more information on configuration and parsing, see https://github.com/willbicks/epigram/blob/main/docs/config.md
package config

import (
	"os"
	"strings"
)

// Repository selects one of a few options for data persistence
type Repository int8

const (
	// Inmemory is a repository that stores data in memory, and is not persistent across restarts
	Inmemory Repository = iota + 1
	// SQLite is a repository that stores data in a SQLite database
	SQLite
)

// repoFromString converts the name of a repository to a variable with type Repository, matching regardless of
// capitalization, and returning 0 if not valid.
func repoFromString(str string) Repository {
	switch strings.ToLower(str) {
	case "inmemory":
		return Inmemory
	case "sqlite":
		return SQLite
	}
	return 0
}

// OIDCProvider provides configuration to initialize a singular OIDC provider
type OIDCProvider struct {
	Name         string `yaml:"name"`
	IssuerURL    string `yaml:"issuerURL"`
	ClientID     string `yaml:"clientID"`
	ClientSecret string `yaml:"clientSecret"`
}

// EntryQuestion is a question the user must answer before being granted entrance to the applicaiton
type EntryQuestion struct {
	Question string `yaml:"question"`
	Answer   string `yaml:"answer"`
}

// Application represents the root configuration struct for the server.
type Application struct {
	// Address is an IP address (or hostname) to bind the server to.
	Address string
	// Port is the port to be bound to on the specified address.
	Port uint16
	// BaseURL is the complete domain and path to access the root of the web server, used for creating callback URLs.
	BaseURL string `yaml:"baseURL"`
	// Title  is the name of the applicaiton to be shown in the frontend.
	Title string `yaml:"title"`
	// Description is a short description of the application to be shown in the frontend.
	Description string `yaml:"description"`
	// Repo dictates what type of storage the application should use for data persistence.
	Repo Repository `yaml:"repo"`
	// DBLoc is the location where the database can be found. In the case of an SQLite repository, this is the path to database file.
	DBLoc string `yaml:"DBLoc"`
	// TrustProxy dictates whether the `X-Forwarded-For`` header should be trusted to obtain the client IP, or if the requestor IP shoud be used instead.
	TrustProxy bool `yaml:"trustProxy"`
	// OIDCProvider is the OIDC provider used to authenticate users.
	OIDCProvider OIDCProvider `yaml:"OIDCProvider"`
	// EntryQuestions is an array of questions.
	EntryQuestions []EntryQuestion `yaml:"entryQuestions"`
}

// merge applies all non-nil / non-default values from the provided layer to the base layer, and returns the result.
//
// Because TrustProxy is a bool, and it's default value is false, a false TrustProxy will not overwrite a base true.
func (base Application) merge(layer Application) Application {
	if layer.Address != "" {
		base.Address = layer.Address
	}
	if layer.Port != 0 {
		base.Port = layer.Port
	}
	if layer.BaseURL != "" {
		base.BaseURL = layer.BaseURL
	}
	if layer.Title != "" {
		base.Title = layer.Title
	}
	if layer.Description != "" {
		base.Description = layer.Description
	}
	if layer.Repo != 0 {
		base.Repo = layer.Repo
	}
	if layer.DBLoc != "" {
		base.DBLoc = layer.DBLoc
	}
	if layer.TrustProxy {
		base.TrustProxy = layer.TrustProxy
	}
	if layer.OIDCProvider != (OIDCProvider{}) {
		base.OIDCProvider = layer.OIDCProvider
	}
	if len(layer.EntryQuestions) > 0 {
		base.EntryQuestions = layer.EntryQuestions
	}
	return base
}

// Parse layers three config sources to return the final application configuration. First, the default configuration is
// loaded (which varied based on system operating system). Then, a .yml configuration file is loaded. If the EP_CONFIG
// env var is set, the yml file is loaded from there, otherwise, it is loaded from the default config location (again,
// varies based on os, see default files for more). Finally, any configuration parameters specified in environment
// variables are applied, and the resulting Application config is returned.
func Parse() (Application, error) {
	// Load default configuration
	cfg := Default

	// Merge YAML configuration
	ymlBytes, err := os.ReadFile(getConfigLoc())
	if err != nil {
		return Application{}, err
	}

	layer, err := parseYAML(ymlBytes)
	if err != nil {
		return Application{}, err
	}
	cfg = cfg.merge(layer)

	// Merge environment configuration
	cfg = cfg.merge(fromEnvironment())

	return cfg, nil
}
