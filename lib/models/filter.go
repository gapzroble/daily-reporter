package models

import "encoding/json"

// Filter report
type Filter struct {
	From            string   `json:"from"`
	To              string   `json:"to"`
	TaskID          []string `json:"taskId"`
	ProjectKey      []string `json:"projectKey"`
	AccountID       []string `json:"accountId"`
	RoleID          []string `json:"roleId"`
	CategoryTypeID  []string `json:"categoryTypeId"`
	FilterID        []string `json:"filterId"`
	ProjectID       []string `json:"projectId"`
	IncludeSubtasks bool     `json:"includeSubtasks"`
	TeamID          []string `json:"teamId"`
	WorkerID        []string `json:"workerId"`
	CustomerID      []string `json:"customerId"`
	CategoryID      []string `json:"categoryId"`
	EpicKey         []string `json:"epicKey"`
	TaskKey         []string `json:"taskKey"`
}

// NewFilter by date and workerID
func NewFilter(date, workerID string) Filter {
	return Filter{
		From:            date,
		To:              date,
		TaskID:          []string{},
		ProjectKey:      []string{},
		AccountID:       []string{},
		RoleID:          []string{},
		CategoryTypeID:  []string{},
		FilterID:        []string{},
		ProjectID:       []string{},
		IncludeSubtasks: false,
		TeamID:          []string{},
		WorkerID:        []string{workerID},
		CustomerID:      []string{},
		CategoryID:      []string{},
		EpicKey:         []string{},
		TaskKey:         []string{},
	}
}

// ToJSONData method
func (filter Filter) ToJSONData() ([]byte, error) {
	data, err := json.Marshal(filter)
	if err != nil {
		return nil, err
	}

	return data, nil
}
