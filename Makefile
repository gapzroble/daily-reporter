SHELL := /bin/bash

build: deps clean test
	#GOARCH=amd64 GOOS=linux go build -o ./bin/DailyReporter ./DailyReporter
	GOARCH=amd64 GOOS=linux go build -o ./bin/Tempo2Nova ./Tempo2Nova
	GOARCH=amd64 GOOS=linux go build -o ./bin/Nova2Csv ./Nova2Csv

run:
	./run.sh

deps:
	GOPRIVATE=github.com go mod vendor

clean:
	ls -I*.sh ./bin | xargs -I {} rm -f ./bin/{}

test:
	# @go test -v $$(go list ./...) >/tmp/gotesting || (grep -A 1 "FAIL:" /tmp/gotesting  && false)
	# @echo PASS
