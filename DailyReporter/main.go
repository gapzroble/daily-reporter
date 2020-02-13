package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os/user"
	"sync"
	"time"
)

func isWeekend(ts time.Time) bool {
	wd := ts.Weekday()
	return wd == time.Sunday || wd == time.Saturday
}

func canRun(ts time.Time) (bool, string) {
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

func main() {
	if ok, reason := canRun(time.Now()); !ok {
		log.Printf("Won't run today: %s\n", reason)
		return
	}

	bufSize := 2
	hours := make(chan float64, bufSize)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		remaining := getLogableHours()
		for i := 0; i < bufSize; i++ {
			hours <- remaining
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		screenshot, err := logNova(hours)
		if err != nil {
			log.Printf("Nova error: %s\n", err.Error())
			return
		}
		usr, _ := user.Current()
		dest := fmt.Sprintf("%s/Downloads/nova_autolog_%s.png", usr.HomeDir, time.Now().Format(time.RFC3339))
		if err := ioutil.WriteFile(dest, screenshot, 0644); err != nil {
			log.Printf("Nova save screenshot: %s\n", err.Error())
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		err := logJira(hours)
		if err != nil {
			log.Printf("JIRA error: %s\n", err.Error())
		}
	}()

	wg.Wait()
}

func handlePanic() {
	msg := recover()
	if msg != nil {
		log.Printf("Go panic: %#v\n", msg)
	}
}
