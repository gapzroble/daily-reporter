package models

import "encoding/json"

// Worklog struct
type Worklog struct {
	IssueKey         string `json:"issueKey"`
	TimeSpentSeconds int64  `json:"timeSpentSeconds"`
	StartDate        string `json:"startDate"`
	StartTime        string `json:"startTime"`
	Description      string `json:"description"`
	AuthorAccountID  string `json:"authorAccountId"`
}

// ToJSONData encode to json
func (w Worklog) ToJSONData() ([]byte, error) {
	dat, err := json.Marshal(w)
	if err != nil {
		return nil, err
	}

	return dat, nil
}

// NewWorklog func
func NewWorklog(date, details, jiraUser string, hours float64) Worklog {
	issue := "TIQ-684" // TODO: map json via env?
	return Worklog{
		IssueKey:         issue,
		TimeSpentSeconds: int64(hours * 60 * 60),
		StartDate:        date,
		StartTime:        "13:00:00",
		Description:      details,
		AuthorAccountID:  jiraUser,
	}
}

// Worklogs struct
type Worklogs struct {
	Results []Worklog `json:"results"`
}

// Logged method
func (logs *Worklogs) Logged() float64 {
	logged := 0.0
	toHours := float64(60 * 60)

	for _, worklog := range logs.Results {
		logged += float64(worklog.TimeSpentSeconds) / toHours
	}

	return logged
}

// NewWorklogs func
func NewWorklogs(data []byte) (*Worklogs, error) {
	logs := &Worklogs{}

	if err := json.Unmarshal(data, logs); err != nil {
		return nil, err
	}

	return logs, nil
}
