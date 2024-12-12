package config

import (
	"os"
	"path"
	"strconv"
	"strings"
)

// EnvironmentPrefix is a string that is prefixed to environment variables seperated by an underscore.
const EnvironmentPrefix string = "EP"

// getConfigLoc returns the location of the config file. If the Config environment variable is set, it will be used.
// Otherwise, the default location will be configDir/config.yml.
func getConfigLoc() string {
	if p := getEnvVar("Config"); p != "" {
		return p
	}
	return path.Join(configDir, "config.yml")
}

// getEnvVar returns the environment variable labeled with the specified name (converted to uppercase), using the
// specified EnvironmentPrefix.
func getEnvVar(name string) string {
	res, _ := os.LookupEnv(strings.Join([]string{EnvironmentPrefix, strings.ToUpper(name)}, "_"))
	return res
}

// fromEnvironment parses some application configuration parameters from similarly named environment variables.
//
// If the value of EnvironmentPrefix is XYZ, the parameter Address will be parsed from the environment variable
// XYZ_ADDRESS.
func fromEnvironment() Application {
	p64, _ := strconv.ParseUint(getEnvVar("Port"), 10, 16)
	port := uint16(p64)

	trustProxy, _ := strconv.ParseBool(getEnvVar("TrustProxy"))
	logJSON, _ := strconv.ParseBool(getEnvVar("LogJSON"))
	devMode, _ := strconv.ParseBool(getEnvVar("DevMode"))

	return Application{
		Title:       getEnvVar("Title"),
		Description: getEnvVar("Description"),
		Address:     getEnvVar("Address"),
		Port:        port,
		BaseURL:     getEnvVar("BaseURL"),
		Repo:        repoFromString(getEnvVar("Repo")),
		DBLoc:       getEnvVar("DBLoc"),
		TrustProxy:  trustProxy,
		LogJSON:     logJSON,
		DevMode:     devMode,
	}
}
