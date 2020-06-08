package schedule

import (
	"testing"
	"time"
)

// TestCanRun test
func TestCanRun(t *testing.T) {
	tests := map[string]bool{
		"2020-01-30": true,
		"2020-01-31": false, // last day
		"2020-02-01": false, // saturday
		"2020-02-02": false, // sunday
		"2020-02-03": true,
		"2020-02-28": false, // last day
		"2020-02-29": false, // saturday
	}
	for date, expected := range tests {
		ts, err := time.Parse("2006-01-02", date)
		if err != nil {
			t.Error(err.Error())
			continue
		}
		actual, _ := CanRun(ts)
		if actual != expected {
			t.Errorf("Expecting %#v on %s, got %#v", expected, date, actual)
		}
	}
}
