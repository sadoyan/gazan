#!/bin/bash

export GOPATH=`pwd`
go get gopkg.in/ini.v1
go get github.com/golang-jwt/jwt
go build -o /tmp/Gazan src/start.go

reflex -r '\.go' -s -- sh -c  'go run src/start.go -up=config/upstreams.json'
