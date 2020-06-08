package tempo

import (
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/rroble/daily-reporter/lib/log"
	"github.com/rroble/daily-reporter/lib/models"
	"github.com/rroble/daily-reporter/lib/schedule"
)

// GetLoggable in seconds
func GetLoggable() (loggable int64, err error) {
	log.Debug("jira", "Getting available hours")
	defer func() {
		log.Debug("jira", "Loggable hours: %.2f", float64(loggable)/3600)
	}()
	loggable = 30600 // default
	r := strings.NewReplacer("{jiraUser}", jiraUser, "{date}", schedule.Today)
	url := r.Replace(getWorklogsURL)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return loggable, err
	}
	req.Header.Add("Authorization", "Bearer "+tempoToken)
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return loggable, err
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return loggable, err
	}

	worklogs, err := models.NewWorklogs(body)
	if err != nil {
		return loggable, err
	}

	return loggable - worklogs.Logged(), nil
}

// Logs all worklogs for the day
func Logs() (logs *models.Worklogs, err error) {
	log.Debug("jira", "Getting worklogs")

	r := strings.NewReplacer("{jiraUser}", jiraUser, "{date}", schedule.Today)
	url := r.Replace(getWorklogsURL)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+tempoToken)
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	worklogs, err := models.NewWorklogs(body)
	if err != nil {
		return nil, err
	}

	return worklogs, nil
}
