package schedule

import (
	"os"
	"time"
)

func isWeekend(ts time.Time) bool {
	wd := ts.Weekday()
	return wd == time.Sunday || wd == time.Saturday
}

// CanRun check if can run now
func CanRun(ts time.Time) (bool, string) {
	// don't check if we specify DATE
	if os.Getenv("DATE") != "" {
		return true, "Date specified"
	}

	// check weekends
	if isWeekend(ts) {
		return false, "Weekend"
	}

	// check last day of the month
	// don't run on schedule(usually 10:30pm)
	// probably logged already
	currentMonth := ts.Month()
	for {
		ts = ts.AddDate(0, 0, 1)
		if isWeekend(ts) {
			continue
		}
		if ts.Month() != currentMonth {
			return false, "Last day of month"
		}
		break // fine
	}

	// TODO: check holidays
	return true, ""
}
