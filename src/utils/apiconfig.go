package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
)

type Ddconf struct {
	Upstreams map[string][]string
	Constants map[string][]string
	sync.RWMutex
}

var Dconf = &Ddconf{
	Upstreams: nil,
	Constants: nil,
	RWMutex:   sync.RWMutex{},
}
var Serob = make(map[string][]string)

func logprint(msg string, err error) {
	if err != nil {
		log.Println(msg, err)
	}
}

func ApiConfig(r *http.Request) {
	Dconf.Lock()
	decoder := json.NewDecoder(r.Body)
	er := decoder.Decode(&Dconf.Upstreams)
	Dconf.Unlock()

	logprint("Json Decode error:", er)
	for k, v := range Dconf.Upstreams {
		Serob[k] = v
		for vv := range v {
			fmt.Println("Registering URL", k, "To Upstream:", v[vv])
		}
		fmt.Println(" ")
	}

}

func valod() {
	fmt.Println("Serob:", Serob)
	fmt.Println("Dconf:", Dconf.Upstreams)
	fmt.Println("----------")
}
