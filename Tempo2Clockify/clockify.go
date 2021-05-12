package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/rroble/daily-reporter/lib/log"
	"github.com/rroble/daily-reporter/lib/models"
	"github.com/rroble/daily-reporter/lib/schedule"
)

var (
	apiKey     string
	simulation bool
	dt         time.Time
	endpoint   = "https://api.clockify.me/api/v1/workspaces/5e37ce129fdab913d809e55f/time-entries"
)

func init() {
	apiKey = os.Getenv("API_KEY")
	if val := os.Getenv("SIMULATION"); val != "" {
		simulation = true
	}
}

// LogFromTempo to clockify
func newTimeEntry(worklog models.Worklog) error {
	if dt.IsZero() {
		dt = schedule.Date
	}

	// break
	if dt.Hour() >= 16 {
		dt = dt.Add(1 * time.Hour)
	}

	start := dt.UTC()
	end := dt.Add(time.Duration(worklog.TimeSpentSeconds) * time.Second).UTC()
	dt = end.Add(1 * time.Second)

	newEntry := &timeEntry{
		Billable:    true,
		Description: worklog.Description,
		Start:       &start,
		End:         &end,
	}

	switch worklog.Issue.Key {
	case "TIQ-684": // Byggmax integration
		newEntry.ProjectID = "60755454fd1fb82c66def7bc"
		newEntry.TaskID = "6075547d8917ba2ef97d3597"

	case "TIQ-957": // Cleveron
		newEntry.ProjectID = "60743a766351ec754bca92d6"
		newEntry.TaskID = "60743a876351ec754bca939a"
		newEntry.Description = "Cleveron: " + worklog.Description
	case "TIQ-1351": // NCP
		newEntry.ProjectID = "60743a766351ec754bca92d6"
		newEntry.TaskID = "60743a876351ec754bca939a"
		newEntry.Description = "NCP OST: " + worklog.Description

	case "TIQ-621": // sysops
	case "TIQ-1095": // training
	case "TIQ-721": // SRS
	case "TIQ-1075": // Sonat change requests
	case "TIQ-?": // Byggmax invoice
	case "TIQ-705": // nk
	case "TIQ-1493": // Accumbo
	case "TIQ-1589": // Arrow
	}

	if err := newEntry.Valid(); err != "" {
		return fmt.Errorf("%s: %s", err, newEntry)
	}

	if simulation {
		log.Debug("clockify", "Simulate log: %s", newEntry)
		return nil
	}

	log.Debug("clockify", "New time entry: %s", newEntry)

	json, err := newEntry.ToJSONData()
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(json))
	if err != nil {
		return fmt.Errorf("New request failed: %s", err.Error())
	}
	req.Header.Add("X-Api-Key", apiKey)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("Do() failed: %s", err.Error())
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("ReadAll() failed: %s", err.Error())
	}

	log.Debug("clockify", "%s", body)

	return nil
}
