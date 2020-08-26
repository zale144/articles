#!/usr/bin/env bash
set -e

echo "#################### downloading CompileDaemon"
GO111MODULE=off go get github.com/githubnemo/CompileDaemon

echo "#################### starting deamon"
CompileDaemon --build="go build -o main cmd/usertags/main.go" --command=./main
