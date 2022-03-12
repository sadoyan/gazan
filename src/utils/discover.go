package utils

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

var seeds = []string{"https://netangels.net/utils/testurl"}
var seed = "https://netangels.net/utils/testurl"
var hostlist []string

type hl struct {
	Hostlist []string
	sync.RWMutex
}

var hoho = hl{
	Hostlist: nil,
	RWMutex:  sync.RWMutex{},
}

func testEq(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func GetSeed() {
	for {
		thishosts := GetHostsByHTTP(seed)
		if !testEq(thishosts, hostlist) {
			hostlist = thishosts
			hoho.RWMutex.Lock()
			hoho.Hostlist = thishosts
			hoho.RWMutex.Unlock()
			fmt.Println("New hosts list is", hoho.Hostlist)
		}
		time.Sleep(10 * time.Second)
	}

}

func RetRandom() string {
	randomIndex := rand.Intn(len(hoho.Hostlist))
	pick := hoho.Hostlist[randomIndex]
	return pick
}
