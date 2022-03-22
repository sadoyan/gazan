#!/bin/bash


case "$1" in

a)
  curl -XPOST -u 'test:Te$ting' --data-binary @config/upstreams2.json 127.0.0.1:4141/config?cfg=append
;;
b)
  curl -XPOST -u 'test:Te$ting' -d {"hololo":"mololo"} http://127.0.0.1:8080/mukuch
  curl -XPOST -u 'test:Te$ting' -d {"vovovo":"cococo"} http://pastor:8080/pooz
  curl -XPOST -u 'test:Te$ting' -d {"jojojo":"xoxoxo"} http://192.168.10.10:8080/jivan
  curl -XPOST -u 'test:Te$ting' -d {"ahahah":"hohoho"} http://192.168.10.10:8080/valer
;;

c)
$0 a
$0 b
;;
*)
echo "Usage: `basename $0` a | b | c"
;;

esac
