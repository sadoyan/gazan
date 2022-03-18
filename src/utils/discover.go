package utils

import (
	"errors"
	"math/rand"
)

//var seeds = []string{"https://netangels.net/utils/testurl"}
//var seed = "https://netangels.net/utils/testurl"
//var hostlist []string
//type hl struct {
//	Hostlist []string
//	sync.RWMutex
//}
//var hoho = hl{
//	Hostlist: nil,
//	RWMutex:  sync.RWMutex{},
//}

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

//func GetSeed() {
//	for {
//		thishosts := GetHostsByHTTP(seed)
//		if !testEq(thishosts, hostlist) {
//			hostlist = thishosts
//			hoho.RWMutex.Lock()
//			hoho.Hostlist = thishosts
//			hoho.RWMutex.Unlock()
//			fmt.Println("New hosts list is", hoho.Hostlist)
//		}
//		time.Sleep(10 * time.Second)
//	}
//
//}

func RetRandomMap(key string) (string, error) {

	if len(Dconf.Upstreams[key]) >= 1 {
		randomIndex := rand.Intn(len(Dconf.Upstreams[key]))
		pick := Dconf.Upstreams[key][randomIndex]
		return pick, nil

	} else {
		return "", errors.New("upstream not found, or upstreams list is empty")
	}

}

func RetRandom(hoho []string) string {
	randomIndex := rand.Intn(len(hoho))
	pick := hoho[randomIndex]
	return pick
}
