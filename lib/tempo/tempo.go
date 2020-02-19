package tempo

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/rroble/daily-reporter/lib/log"
	"github.com/rroble/daily-reporter/lib/models"
)

// Log tempo
func Log(hours <-chan float64) error {
	defer log.Debug("tempo", "Done Tempo")
	logHours := <-hours
	if logHours <= 0 {
		return nil
	}
	log.Debug("tempo", "Logging Tempo")
	data, err := models.NewWorklog(today, details, jiraUser, logHours).ToJSONData()
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
