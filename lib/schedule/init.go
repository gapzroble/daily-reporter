package schedule

import (
	"os"
	"strings"
	"time"
)

// Date run date
var Date time.Time

func init() {
	date := os.Getenv("DATE")
	switch strings.ToLower(date) {
	case "", "today":
		Date = time.Now().Local()
	case "yesterday":
		Date = time.Now().Local().AddDate(0, 0, -1)
	case "tomorrow":
		Date = time.Now().Local().AddDate(0, 0, 1)
	default:
		Date, _ = time.Parse("2006-01-02", date)
	}

	if Date.IsZero() {
		panic("Invalid date")
	}

	// shift starts at 13:00 PST
	Date = time.Date(Date.Year(), Date.Month(), Date.Day(), 13, 0, 0, 0, time.Local)
}

// Today but actually the run date
func Today() string {
	return Date.Format("2006-01-02")
}

// DateTime func
func DateTime() string {
	return Date.Format(time.RFC3339)
}
