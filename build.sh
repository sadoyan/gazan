#!/bin/bash

export GOPATH=`pwd`

#go get github.com/go-sql-driver/mysql
#go get gopkg.in/gomail.v2
#go get github.com/google/uuid
go get gopkg.in/ini.v1

go build -o /tmp/GOProxy src/start.go

