package schedule

import (
	"os"
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
	if Today == "" {
		Today = time.Now().Format("2006-01-02")
	}
	Now = Today + time.Now().Format("T15:04:05Z07:00")
}
