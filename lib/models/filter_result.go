package models

import (
	"encoding/json"
	"io"
	"io/ioutil"
)

// FilterResult struct
type FilterResult struct {
	FilterKey string `json:"filterKey"`
}

// NewFilterResult func
func NewFilterResult(r io.Reader) (result FilterResult) {
	dat, err := ioutil.ReadAll(r)
	if err != nil {
		return
	}

	json.Unmarshal(dat, &result)

	return result
}
