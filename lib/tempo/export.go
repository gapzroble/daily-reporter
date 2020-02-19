package tempo

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/rroble/daily-reporter/lib/models"
)

func getFilterKey(date, jiraUser, jwt string) (string, error) {
	filter, err := models.NewFilter(date, jiraUser).ToJSONData()
	if err != nil {
		return "", err
	}

	debug("jira", "Filter report..")
	req, err := http.NewRequest("POST", filterURL, bytes.NewBuffer(filter))
	if err != nil {
		return "", err
	}

	req.Header.Add("Authorization", "JWT "+jwt)
	req.Header.Add("Content-Type", "application/json")
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	if response.StatusCode > 299 {
		return "", fmt.Errorf("Expecting 2xx response code, got %d", response.StatusCode)
	}

	return models.NewFilterResult(response.Body).FilterKey, nil
}

func exportReport(filterKey, jwt string) ([]byte, error) {
	r := strings.NewReplacer("{filterKey}", filterKey, "{jwt}", jwt)
	export := r.Replace(exportURL)
	debug("jira", "Exporting report..")

	response, err := http.Get(export)
	if err != nil {
		return nil, err
	}

	if response.StatusCode > 299 {
		return nil, fmt.Errorf("Expecting 2xx response code, got %d", response.StatusCode)
	}

	return ioutil.ReadAll(response.Body)
}
