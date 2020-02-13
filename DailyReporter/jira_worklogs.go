package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type worklogs struct {
	Results []worklog `json:"results"`
}

// Logged method
func (logs *worklogs) Logged() float64 {
	logged := 0.0
	toHours := float64(60 * 60)

	for _, worklog := range logs.Results {
		logged += float64(worklog.TimeSpentSeconds) / toHours
	}

	return logged
}

func getLogableHours() (hours float64, err error) {
	debug("tempo", "Getting available hours")
	defer func() {
		debug("tempo", "Loggable hours: %.2f", hours)
	}()
	hours = 8.5 // default
	today := time.Now().Format("2006-01-02")
	url := fmt.Sprintf("https://api.tempo.io/core/3/worklogs/user/%s?from=%s&to=%s", jiraUser, today, today)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return hours, err
	}
	req.Header.Add("Authorization", "Bearer "+tempoToken)
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return hours, err
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return hours, err
	}

	worklogs, err := newWorklogs(body)
	if err != nil {
		return hours, err
	}

	return hours - worklogs.Logged(), nil
}

func newWorklogs(data []byte) (*worklogs, error) {
	logs := &worklogs{}

	if err := json.Unmarshal(data, logs); err != nil {
		return nil, err
	}

	return logs, nil
}
