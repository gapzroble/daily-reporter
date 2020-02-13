#!/usr/bin/env bash

curl \
    -H "Authorization: Bearer ${TOKEN}" \
    "https://api.tempo.io/core/3/worklogs/user/557058:ddf95d2c-e3e8-4380-9456-d191554f48b7?from=2020-02-06&to=2020-02-07"
