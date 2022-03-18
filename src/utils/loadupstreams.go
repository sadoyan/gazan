package utils

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

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
