package tempo

import (
	"os"
)

var (
	getWorklogsURL = "https://api.tempo.io/core/3/worklogs/user/{jiraUser}?from={date}&to={date}"
	worklogsURL    = "https://api.tempo.io/core/3/worklogs"
	tempoToken     = ""
	details        = "(autolog)"
	jiraUser       = "557058:ddf95d2c-e3e8-4380-9456-d191554f48b7"
	email          = "r.roble@arcanys.com"
	password       = ""
	loginURL       = "https://id.atlassian.com/login?continue="
	reportID       = "f629bffe-eb81-4097-bc4b-c29c2f563090"
	jiraURL        = "https://arcanys.atlassian.net/plugins/servlet/ac/io.tempo.jira/tempo-app#!/reports/logged-time/{reportID}?columns=WORKED_COLUMN&dateDisplayType=days&from={today}&groupBy=project&groupBy=issue&groupBy=worklog&periodType=FIXED&subPeriodType=MONTH&to={today}&viewType=TIMESHEET&workerId={jiraUser}"
	filterURL      = "https://app.tempo.io/rest/tempo-timesheets/4/worklogs/export/filter"
	exportURL      = "https://app.tempo.io/rest/tempo-timesheets/4/worklogs/export/{filterKey}?format=pdf&title=Daily%2520Report&jwt={jwt}&groupBy=project,issue,worklog&columns=worked"
)

func init() {
	tempoToken = os.Getenv("TEMPO_TOKEN")
	password = os.Getenv("JIRA_PASSWORD")
	if val := os.Getenv("JIRA_EMAIL"); val != "" {
		email = val
	}
	if val := os.Getenv("TEMPO_REPORT"); val != "" {
		reportID = val
	}
}
