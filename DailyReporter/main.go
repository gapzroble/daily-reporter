package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		screenshot, err := logNova()
		if err != nil {
			log.Printf("Nova error: %s\n", err.Error())
			return
		}
		dest := fmt.Sprintf("nova_autolog_%s.png", time.Now().Format(time.RFC3339))
		if err := ioutil.WriteFile(dest, screenshot, 0644); err != nil {
			log.Printf("Nova save screenshot: %s\n", err.Error())
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
