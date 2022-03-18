#!/bin/bash
curl -XPOST -u 'test:Te$ting' --data-binary @config/upstreams.json 127.0.0.1:4141/config
curl -XPOST -u 'test:Te$ting' -d {"hololo":"mololo"} http://127.0.0.1:8080/mukuch
curl -XPOST -u 'test:Te$ting' -d {"vovovo":"cococo"} http://pastor:8080/pooz
curl -XPOST -u 'test:Te$ting' -d {"jojojo":"xoxoxo"} http://192.168.10.10:8080/jivan
