package main

import (
	"strings"
	"time"

	"github.com/rroble/daily-reporter/lib/log"
	"github.com/rroble/daily-reporter/lib/nova"
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

	worklogs, err := tempo.Logs()
	if err != nil || worklogs == nil {
		log.Debug("main", "Failed to get worklogs: %s", err.Error())
		return
	}

	nova.Init()
	defer nova.End()

	for _, worklog := range worklogs.Results {
		if strings.HasPrefix(worklog.Issue.Key, "BLOCAL-") {
			continue
		}
		log.Debug("main", "Copy worklog: %s", worklog)
		if err := nova.LogFromTempo(worklog); err != nil {
			log.Debug("main", "Failed to log nova: %s", err.Error())
		}
	}

}

func handlePanic() {
	msg := recover()
	if msg != nil {
		log.Debug("main", "Go panic: %#v", msg)
	}
}
