#!/usr/bin/env bash

GOFMT_OUTPUT="$(gofmt -l . 2>&1)"
if [ -n "$GOFMT_OUTPUT" ]; then
   echo "All the following files are not correctly formatted"
   echo "${GOFMT_OUTPUT}"
   exit 1
fi

