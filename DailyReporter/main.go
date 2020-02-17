package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"time"
)

func main() {
	defer handlePanic()

	if ok, reason := canRun(time.Now()); !ok {
		debug("main", "Won't run today: %s", reason)
		return
	}

	bufSize := 2 // temp and nova threads
	hours := make(chan float64, bufSize)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		remaining, err := getLogableHours()
		if err != nil {
			debug("main", "Worklogs error: %s", err.Error())
			// still continue
		}
		// unless no remaining hours
		// if remaining <= 0 {
		// 	debug("main", "No loggable hours, already logged?")
		// 	debug("main", "Quit")
		// 	os.Exit(-1)
		// }
		for i := 0; i < bufSize; i++ {
			hours <- remaining
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		screenshot, err := logNova(hours)
		if err != nil {
			debug("main", "Failed to log nova, %s", err.Error())
			return
		}
		if screenshot == nil {
			debug("main", "No nova screenshot found")
			return
		}
		dest := fmt.Sprintf("autolog_nova_%s.png", now)
		if err := ioutil.WriteFile(dest, screenshot, 0644); err != nil {
			debug("main", "Failed to save nova screenshot, %s", err.Error())
		}
	}()

	doneTempo := make(chan bool)

	wg.Add(1)
	go func() {
		defer wg.Done()
		err := logTempo(hours)
		if err != nil {
			debug("main", "Failed to log tempo, %s", err.Error())
			doneTempo <- false
			return
		}
		doneTempo <- true
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		screenshot, err := jiraScreenshot(doneTempo)
		if err != nil {
			debug("main", "Failed to create jira screenshot, %s", err.Error())
			return
		}
		dest := fmt.Sprintf("autolog_jira_%s.png", now)
		if err := ioutil.WriteFile(dest, screenshot, 0644); err != nil {
			debug("main", "Failed to save jira screenshot, %s", err.Error())
		}
	}()

	wg.Wait()
}

func debug(thread, msg string, args ...interface{}) {
	log.Printf("["+thread+"] "+msg+"\n", args...)
}

func isWeekend(ts time.Time) bool {
	wd := ts.Weekday()
	return wd == time.Sunday || wd == time.Saturday
}

func canRun(ts time.Time) (bool, string) {
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

func handlePanic() {
	msg := recover()
	if msg != nil {
		debug("main", "Go panic: %#v", msg)
	}
}
