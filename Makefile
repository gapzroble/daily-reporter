SHELL := /bin/bash

build: deps clean test
	GOARCH=amd64 GOOS=linux go build -o ./bin/DailyReporter ./DailyReporter

run:
	go run ./DailyReporter/*.go

deps:
	GOPRIVATE=github.com go mod vendor

clean:
	ls -I*.sh ./bin | xargs -I {} rm -f ./bin/{}

test:
	@go test -v $$(go list ./...) >/tmp/gotesting || (grep -A 1 "FAIL:" /tmp/gotesting  && false)
	@echo PASS
