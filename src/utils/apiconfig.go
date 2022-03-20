package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func ApiConfig(r *http.Request) {
	Dconf.Lock()
	decoder := json.NewDecoder(r.Body)
	er := decoder.Decode(&Dconf.Upstreams)
	Dconf.Unlock()
	if er != nil {
		log.Println("Json Decode error:", er)
	}
	for k, v := range Dconf.Upstreams {
		Serob[k] = v
		for vv := range v {
			fmt.Println("Registering URL", k, "To Upstream:", v[vv])
		}
		fmt.Println(" ")
	}

}

func LoadUpstreams(up string) {
	data, err := ioutil.ReadFile(up)
	if err != nil {
		log.Println("Cant load default upstreams file:", err)
		log.Println("Startingwithout upstreams")
	} else {
		er := json.Unmarshal([]byte(data), &Dconf.Upstreams)
		if er == nil {
			log.Println("Sucesfully loaded default upstrems list")
		} else {
			log.Println("Error decoding default upstreams list")
		}
	}

}

func valod() {
	fmt.Println("Serob:", Serob)
	fmt.Println("Dconf:", Dconf.Upstreams)
	fmt.Println("----------")
}
