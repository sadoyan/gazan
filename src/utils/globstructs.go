package utils

import (
	"sync"
)

type Ddconf struct {
	Upstreams map[string][]string
	Constants map[string][]string
	sync.RWMutex
}

var Dconf = &Ddconf{
	Upstreams: make(map[string][]string),
	Constants: make(map[string][]string),
	RWMutex:   sync.RWMutex{},
}

//var Serob = make(map[string][]string)
