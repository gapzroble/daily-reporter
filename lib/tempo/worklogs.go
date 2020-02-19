package tempo

import (
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/rroble/daily-reporter/lib/log"
	"github.com/rroble/daily-reporter/lib/models"
)

// GetLoggableHours func
func GetLoggableHours() (hours float64, err error) {
	log.Debug("tempo", "Getting available hours")
	defer func() {
		log.Debug("tempo", "Loggable hours: %.2f", hours)
	}()
	hours = 8.5 // default
	r := strings.NewReplacer("{jiraUser}", jiraUser, "{date}", today)
	url := r.Replace(getWorklogsURL)

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

	worklogs, err := models.NewWorklogs(body)
	if err != nil {
		return hours, err
	}

	return hours - worklogs.Logged(), nil
}
