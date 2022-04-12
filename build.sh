#!/bin/bash
export GAZANKEY='a66cd04bb85a2daed5080fb41c3da6642f37f4390d76e37c2a57f4edd4c9324e'
export BASICUSER='test'
export BASICPASS='Te$ting'
export JWTSECRET='Super$ecter123765@'


export GOPATH=`pwd`
go get gopkg.in/ini.v1
go get github.com/golang-jwt/jwt
go build -o /tmp/Gazan src/start.go

#reflex -d none -r '.go'  -v  -s -- sh -c  'go run src/start.go -up=config/upstreams.json'

reflex -d none -r '.'  -s -- sh -c  'go run src/start.go -config config.ini -up=config/upstreams.json'
#reflex -d fancy -r '.'  -s -- sh -c  'go run src/start.go -up=config/upstreams.json'




# curl -XPOST -u 'test:Te$ting' --data-binary @config/upstreams.json 127.0.0.1:4141/config?cfg=new