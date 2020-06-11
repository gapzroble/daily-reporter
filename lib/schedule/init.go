package schedule

import (
	"os"
	"strings"
	"time"
)

var (
	// Today date
	Today string

	// Now date + time
	Now string
)

func init() {
	Today = os.Getenv("DATE")
	if Today == "" || strings.ToLower(Today) == "today" {
		Today = time.Now().Format("2006-01-02")
	}
	if Today == "" || strings.ToLower(Today) == "yesterday" {
		Today = time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	}
	Now = Today + time.Now().Format("T15:04:05Z07:00")
}
