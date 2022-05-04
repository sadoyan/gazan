package utils

import (
	"sync"
)

type Ddconf struct {
	Upstreams map[string][]string
	Constants map[string][]string
	Windcards map[string]bool
	sync.RWMutex
}

var Dconf = &Ddconf{
	Upstreams: make(map[string][]string),
	Constants: make(map[string][]string),
	Windcards: make(map[string]bool),
	RWMutex:   sync.RWMutex{},
}

//var Serob = make(map[string][]string)
