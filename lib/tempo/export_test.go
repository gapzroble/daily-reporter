package tempo

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/rroble/daily-reporter/lib/log"
)

// TestPDF test
func TestPDF(t *testing.T) {
	data, err := ioutil.ReadFile("../../autolog_jira_2020-02-19T16:38:19+08:00.pdf")
	if err != nil {
		t.Errorf("Failed to load pdf file, %s", err.Error())
		return
	}

	screenshot, err := convertReport(data)
	if err != nil {
		t.Errorf("Failed to convert pdf to image, %s", err.Error())
	}
	dest := fmt.Sprintf("../../autolog_jira_%s.png", "test")
	if err := ioutil.WriteFile(dest, screenshot, 0644); err != nil {
		log.Debug("main", "Failed to save jira screenshot, %s", err.Error())
	}
}
