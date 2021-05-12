package main

import (
	"encoding/json"
	"time"
)

type timeEntry struct {
	Start       *time.Time `json:"start"`
	Billable    bool       `json:"billable"`
	Description string     `json:"description"`
	ProjectID   string     `json:"projectId"`
	TaskID      string     `json:"taskId"`
	End         *time.Time `json:"end"`
}

func (e *timeEntry) ToJSONData() ([]byte, error) {
	out, err := json.Marshal(e)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (e *timeEntry) String() string {
	out, _ := e.ToJSONData()
	return string(out)
}

func (e timeEntry) Valid() string {
	if e.ProjectID == "" {
		return "No projectId"
	}
	if e.TaskID == "" {
		return "No taskId"
	}

	if e.Start == nil {
		return "No start"
	}
	if e.End == nil {
		return "No end"
	}
	if e.Start.After(*e.End) {
		return "Start > End"
	}

	return ""
}
