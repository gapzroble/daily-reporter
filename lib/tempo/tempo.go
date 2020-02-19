package tempo

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"

	"github.com/rroble/daily-reporter/lib/log"
	"github.com/rroble/daily-reporter/lib/models"
	"github.com/rroble/daily-reporter/lib/schedule"
)

// Log tempo
func Log(loggable <-chan int64) error {
	defer log.Debug("jira", "Done Tempo")
	logSeconds := <-loggable
	if logSeconds <= 0 {
		return errors.New("Already logged")
	}
	log.Debug("jira", "Logging Tempo")
	data, err := models.NewWorklog(schedule.Today, details, jiraUser, logSeconds).ToJSONData()
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", worklogsURL, bytes.NewBuffer(data))
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
