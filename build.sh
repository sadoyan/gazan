#!/bin/bash

export GOPATH=`pwd`
go get gopkg.in/ini.v1

go build -o /tmp/Gazan src/start.go

reflex -r '\.go' -s -- sh -c  'go run src/start.go'
