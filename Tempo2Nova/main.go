package main

import (
	"fmt"
	"io/ioutil"
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
	log.Debug("main", "Date is %s", schedule.Today)

	worklogs, err := tempo.Logs()
	if err != nil {
		log.Debug("main", "Failed to get worklogs: %s", err.Error())
		return
	}
	if worklogs == nil {
		log.Debug("main", "No worklogs found")
		return
	}

	nova.Init()
	defer nova.End()

	for _, worklog := range worklogs.Results {
		if !strings.HasPrefix(worklog.Issue.Key, "TIQ-") {
			continue
		}

		// skip, already logged
		// if strings.HasPrefix(worklog.Issue.Key, "TIQ-") {
		// 	continue
		// }

		log.Debug("main", "Copy worklog: %s", worklog)
		if err := nova.LogFromTempo(worklog); err != nil {
			log.Debug("main", "Failed to log nova: %s", err.Error())
		}
	}

	screenshot, err := nova.PrintScreen()
	if err != nil {
		log.Debug("main", "Failed to log nova, %s", err.Error())
		return
	}
	if screenshot == nil {
		log.Debug("main", "No nova screenshot found")
		return
	}
	dest := fmt.Sprintf("/home/randolph/Downloads/nova_%s.png", schedule.Now)
	if err := ioutil.WriteFile(dest, screenshot, 0644); err != nil {
		log.Debug("main", "Failed to save nova screenshot, %s", err.Error())
	}
}

func handlePanic() {
	msg := recover()
	if msg != nil {
		log.Debug("main", "Go panic: %#v", msg)
	}
}
