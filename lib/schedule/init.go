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
	switch strings.ToLower(Today) {
	case "", "today":
		Today = time.Now().Format("2006-01-02")
	case "yesterday":
		Today = time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	case "tomorrow":
		Today = time.Now().AddDate(0, 0, 1).Format("2006-01-02")
	}
	Now = Today + time.Now().Format("T15:04:05Z07:00")
}
