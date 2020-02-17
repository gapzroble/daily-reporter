package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type worklog struct {
	IssueKey         string `json:"issueKey"`
	TimeSpentSeconds int64  `json:"timeSpentSeconds"`
	StartDate        string `json:"startDate"`
	StartTime        string `json:"startTime"`
	Description      string `json:"description"`
	AuthorAccountID  string `json:"authorAccountId"`
}

func (w worklog) toJSON() ([]byte, error) {
	dat, err := json.MarshalIndent(w, "", "\t")
	if err != nil {
		return nil, err
	}

	return dat, nil
}

func newWorklog(hours float64) worklog {
	issue := "TIQ-684" // TODO: map json via env?
	return worklog{
		IssueKey:         issue,
		TimeSpentSeconds: int64(hours * 60 * 60),
		StartDate:        today,
		StartTime:        "13:00:00",
		Description:      details,
		AuthorAccountID:  jiraUser,
	}
}

func logTempo(hours <-chan float64) error {
	defer debug("tempo", "Done Tempo")
	logHours := <-hours
	if logHours <= 0 {
		return nil
	}
	debug("tempo", "Logging Tempo")
	data, err := newWorklog(logHours).toJSON()
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", "https://api.tempo.io/core/3/worklogs", bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+tempoToken)
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if response.StatusCode > 299 {
		return fmt.Errorf("Expecting 2xx response code, got %d", response.StatusCode)
	}

	return nil
}
