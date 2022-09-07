#!/bin/bash

#export GAZANKEY='a66cd04bb85a2daed5080fb41c3da6642f37f4390d76e37c2a57f4edd4c9324e'
#export BASICUSER='test'
#export BASICPASS='Te$ting'
#export JWTSECRET='Super$ecter123765@'
export CONFIGKEY='47690a568c204f1b85fe997ce9fc3ee53a571be5508a4d8ca1da238533ffe3f2'
#
#
#export GOPATH=/opt/GOLang/goext
#go build -o /tmp/Gazan main.go
#
#reflex -d none -r '.'  -s -- sh -c  'go run ./ -config config.ini -up=cfgjson/upstreams.json'
reflex -d none -r '.'  -s -- sh -c  'go run ./ -config config.yml -up=cfgjson/upstreams.json'
#reflex -d none -r '.'  -s -- sh -c  'rsync -vzal ../gazan razor:/usr/local/src/ --exclude config.ini --exclude build.sh  && sleep 10d'

# curl -XPOST -u 'test:Te$ting' --data-binary @config/upstreams.json 127.0.0.1:4141/config?cfg=new
