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
	sync.RWMutex
}

var Dconf = Ddconf{
	Upstreams: nil,
	RWMutex:   sync.RWMutex{},
}

func logprint(msg string, err error) {
	if err != nil {
		log.Println(msg, err)
	}
}

func ApiConfig(r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	er := decoder.Decode(&Dconf.Upstreams)
	logprint("Json Decode error:", er)

	for k, v := range Dconf.Upstreams {
		for vv := range v {
			fmt.Println("Registering URL", k, "To Upstream:", v[vv])
		}
		fmt.Println(" ")
	}

}
