package main

import (
	"fmt"
	"io/ioutil"
	"sync"
	"time"

	"github.com/rroble/daily-reporter/lib/log"
	"github.com/rroble/daily-reporter/lib/nova"
	"github.com/rroble/daily-reporter/lib/tempo"
)

func main() {
	defer handlePanic()

	if ok, reason := canRun(time.Now()); !ok {
		log.Debug("main", "Won't run today: %s", reason)
		return
	}

	bufSize := 2 // temp and nova threads
	hours := make(chan float64, bufSize)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		// ------------ test -------------
		// var remaining = 0.0
		// var err error
		// ------------ test -------------
		remaining, err := tempo.GetLoggableHours()
		if err != nil {
			log.Debug("main", "Worklogs error: %s", err.Error())
			// still continue
		}
		for i := 0; i < bufSize; i++ {
			hours <- remaining
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		screenshot, err := nova.Log(hours)
		if err != nil {
			log.Debug("main", "Failed to log nova, %s", err.Error())
			return
		}
		if screenshot == nil {
			log.Debug("main", "No nova screenshot found")
			return
		}
		dest := fmt.Sprintf("autolog_nova_%s.png", now)
		if err := ioutil.WriteFile(dest, screenshot, 0644); err != nil {
			log.Debug("main", "Failed to save nova screenshot, %s", err.Error())
		}
	}()

	doneTempo := make(chan bool)

	wg.Add(1)
	go func() {
		defer wg.Done()
		err := tempo.Log(hours)
		if err != nil {
			log.Debug("main", "Failed to log tempo, %s", err.Error())
			doneTempo <- false
			return
		}
		doneTempo <- true
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		screenshot, err := tempo.Report(doneTempo)
		if err != nil {
			log.Debug("main", "Failed to create jira screenshot, %s", err.Error())
			return
		}
		dest := fmt.Sprintf("autolog_jira_%s.png", now)
		if err := ioutil.WriteFile(dest, screenshot, 0644); err != nil {
			log.Debug("main", "Failed to save jira screenshot, %s", err.Error())
		}
	}()

	wg.Wait()
}

func handlePanic() {
	msg := recover()
	if msg != nil {
		log.Debug("main", "Go panic: %#v", msg)
	}
}
