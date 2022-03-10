package config

import (
	"os"
	"strings"
)

const EnvironmentPrefix string = "EP"

func fromEnvironment() Application {
	// getEnvVar returns the environment variable labeled with the specified name (converted to uppercase), using the
	// specified EnvironmentPrefix.
	getEnvVar := func(name string) interface{} {
		res, _ := os.LookupEnv(strings.Join([]string{EnvironmentPrefix, strings.ToUpper(name)}, "_"))
		return res
	}

	return Application{
		Address:    getEnvVar("Address").(string),
		Port:       getEnvVar("Port").(uint16),
		BaseURL:    getEnvVar("BaseURL").(string),
		Repo:       repoFromString(getEnvVar("Repo").(string)),
		DBLoc:      getEnvVar("DBLoc").(string),
		TrustProxy: getEnvVar("TrustProxy").(bool),
	}
}
