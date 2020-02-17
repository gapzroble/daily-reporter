package main

import (
	"os"
	"time"
)

var (
	username     string
	password     string
	project      string
	details      string
	jiraUser     string
	tempoToken   string
	today        string
	now          string
	email        string
	jiraPassword string
)

var (
	width  int64 = 1390
	height int64 = 895
)

func init() {
	username = env("NOVA_USERNAME", "Randolph Roble")
	password = env("NOVA_PASSWORD", "")
	project = env("PROJECT", "Byggmax")
	details = env("DETAILS", "(autolog)")
	jiraUser = env("JIRA_USER", "557058:ddf95d2c-e3e8-4380-9456-d191554f48b7")
	tempoToken = env("TEMPO_TOKEN", "")
	today = env("DATE", time.Now().Format("2006-01-02"))
	now = today + time.Now().Format("T15:04:05Z07:00")
	email = env("JIRA_EMAIL", "r.roble@arcanys.com")
	jiraPassword = env("JIRA_PASSWORD", "")
}

func env(envKey, defaultValue string) string {
	val := os.Getenv(envKey)
	if val == "" {
		return defaultValue
	}

	return val
}
