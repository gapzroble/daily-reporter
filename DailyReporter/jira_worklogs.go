package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
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

func getLogableHours() (hours float64) {
	log.Println("Getting available hours")
	defer func() {
		log.Println("Loggable hours", hours)
	}()
	hours = 8.5 // default
	today := time.Now().Format("2006-01-02")
	url := fmt.Sprintf("https://api.tempo.io/core/3/worklogs/user/%s?from=%s&to=%s", jiraUser, today, today)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("Error creating request, %s", err.Error())
		return
	}
	req.Header.Add("Authorization", "Bearer "+tempoToken)
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("Error creating request, %s", err.Error())
		return
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Printf("Error reading response, %s", err.Error())
		return
	}

	worklogs, err := newWorklogs(body)
	if err != nil {
		log.Printf("Error unmarshall worklogs, %s", err.Error())
		return
	}

	return hours - worklogs.Logged()
}

func newWorklogs(data []byte) (*worklogs, error) {
	logs := &worklogs{}

	if err := json.Unmarshal(data, logs); err != nil {
		return nil, err
	}

	return logs, nil
}
