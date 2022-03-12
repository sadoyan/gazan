#!/bin/bash

export GOPATH=`pwd`
go get gopkg.in/ini.v1

go build -o /tmp/Gazan src/start.go

