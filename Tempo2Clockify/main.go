package main

import (
	"strings"
	"time"

	"github.com/rroble/daily-reporter/lib/log"
	"github.com/rroble/daily-reporter/lib/schedule"
	"github.com/rroble/daily-reporter/lib/tempo"
)

func main() {
	defer handlePanic()
	defer log.Debug("main", "Done")

	if ok, reason := schedule.CanRun(time.Now()); !ok {
		log.Debug("main", "Won't run today: %s", reason)
		return
	}
	log.Debug("main", "Date is %s", schedule.Today())

	worklogs, err := tempo.Logs()
	if err != nil {
		log.Debug("main", "Failed to get worklogs: %s", err.Error())
		return
	}
	if worklogs == nil {
		log.Debug("main", "No worklogs found")
		return
	}

	for _, worklog := range worklogs.Results {
		if !strings.HasPrefix(worklog.Issue.Key, "TIQ-") {
			continue
		}

		// skip, already logged
		// if strings.HasPrefix(worklog.Issue.Key, "TIQ-") {
		// 	continue
		// }

		log.Debug("main", "Copy worklog: %s", worklog)
		if err := newTimeEntry(worklog); err != nil {
			log.Debug("main", "Failed to add time entry: %s", err.Error())
		}
	}

}

func handlePanic() {
	msg := recover()
	if msg != nil {
		log.Debug("main", "Go panic: %#v", msg)
	}
}
