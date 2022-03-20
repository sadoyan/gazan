package utils

import (
	"sync"
)

type Ddconf struct {
	Upstreams map[string][]string
	//Constants map[string][]string
	sync.RWMutex
}

var Dconf = &Ddconf{
	Upstreams: nil,
	//Constants: nil,
	RWMutex: sync.RWMutex{},
}
var Serob = make(map[string][]string)
