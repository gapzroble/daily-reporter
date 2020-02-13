package main

import "os"

var (
	username   string
	password   string
	project    string
	details    string
	jiraUser   string
	tempoToken string
)

func init() {
	username = env("NOVA_USERNAME", "Randolph Roble")
	password = env("NOVA_PASSWORD", "")
	project = env("PROJECT", "byggmax")
	details = env("DETAILS", "(autolog)")
	jiraUser = env("JIRA_USER", "557058:ddf95d2c-e3e8-4380-9456-d191554f48b7")
	tempoToken = env("TEMPO_TOKEN", "")
}

func env(envKey, defaultValue string) string {
	val := os.Getenv(envKey)
	if val == "" {
		return defaultValue
	}

	return val
}
