#!/usr/bin/env bash
set -e

echo "#################### running tests"
NF_DB_NAME=test NF_CONFIG=../config/config.yaml CGO_ENABLED=0 go test -p 1 -count=1 -tags=integration ./... -v --cover

echo "#################### downloading CompileDaemon"
GO111MODULE=off go get github.com/githubnemo/CompileDaemon

echo "#################### starting deamon"
CompileDaemon --build="go build -o main cmd/newsfeed/main.go" --command=./main
