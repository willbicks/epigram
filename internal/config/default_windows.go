//go:build windows

package config

// configLoc is the default location to search for config files in
const configLoc = "."

// Default is a default configuration, used as a base for additional configurations to be merged on top of.
var Default = Application{
	Address:     "0.0.0.0",
	Port:        80,
	Title:       "Epigram",
	Description: "Epigram is a simple web service for communities to immortalize the enlightening, funny, or downright dumb quotes that they hear.",
	TrustProxy:  false,
	Repo:        SQLite,
	DBLoc:       ".",
}
