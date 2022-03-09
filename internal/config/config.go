// package config is responsible for processing global, application configuration determined at runtime.
// Supports both reading from .yml files, and environment variables, the latter taking priority over the former.
package config

// Repository selects one of a few options for data persistance
type Repository int8

const (
	Inmemory Repository = iota + 1
	SQLite
)

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
	BaseURL        string          `yaml:"baseURL"`
	Title          string          `yaml:"title"`
	Description    string          `yaml:"description"`
	Repo           Repository      `yaml:"repo"`
	TrustProxy     bool            `yaml:"trustProxy"`
	OIDCProvider   OIDCProvider    `yaml:"OIDCProvider"`
	EntryQuestions []EntryQuestion `yaml:"entryQuestions"`
}
