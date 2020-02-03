#!/usr/bin/env bash

curl -i \
    -X POST \
    -H "Authorization: Bearer ${TOKEN}" \
    -H "Content-Type: application/json" \
    -d @worklog.json \
    https://api.tempo.io/core/3/worklogs
