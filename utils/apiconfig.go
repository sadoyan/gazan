package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"reflect"
	"strings"
	"time"
)

func ApiConfig(r *http.Request) {
	urlparam := r.URL.Query().Get("cfg")
	decoder := json.NewDecoder(r.Body)

	switch urlparam {
	case "append":
		tempUps := make(map[string][]string)
		er := decoder.Decode(&tempUps)
		if er != nil {
			log.Println("Json Decode error:", er)
		}
		var changes bool
		for k, v := range tempUps {
			for vv := range v {
				switch Contains(Dconf.Constants[k], v[vv]) {
				case false:
					Dconf.Lock()
					Dconf.Constants[k] = append(tempUps[k], v[vv])
					Dconf.Upstreams[k] = append(tempUps[k], v[vv])
					Dconf.Unlock()
					fmt.Println("Registering URL", k, "To Upstream:", v[vv])
					changes = true
				}
			}
		}
		if !changes {
			log.Println("No Changes sice last update")
		}
	case "new":
		tempUps := make(map[string][]string)
		er := decoder.Decode(&tempUps)
		if er != nil {
			fmt.Println("Error decoding json:", er)
		}
		fmt.Println("Creating new upstream config")
		Dconf.Lock()
		Dconf.Constants = tempUps
		Dconf.Upstreams = tempUps
		Dconf.Unlock()
		for k, x := range Dconf.Constants {
			log.Println("Main: ", k)
			for v := range x {
				log.Println("  Upstream: ", x[v])
			}
		}
	default:
		log.Println("Unknown parameter ")
	}

}

func LoadUpstreamsFronFIle(up string) {
	data, err := ioutil.ReadFile(up)
	if err != nil {
		log.Println("Cant load default upstreams file:", err)
		log.Println("Startingwithout upstreams")
	} else {
		er := json.Unmarshal(data, &Dconf.Upstreams)
		for k, v := range Dconf.Upstreams {
			if !strings.Contains(k, "/") {
				Dconf.Windcards[k] = true
			}
			Dconf.Constants[k] = v
			for vv := range v {
				fmt.Println("Registering URL", k, "To Upstream:", v[vv])
			}
			fmt.Println(" ")
		}

		if er == nil {
			log.Println("Sucesfully loaded default upstrems list")
		} else {
			log.Println("Error decoding default upstreams list")
		}
	}
}

func Valod(healtchecks int) {
	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   2 * time.Second,
			KeepAlive: 2 * time.Second,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		MaxIdleConnsPerHost:   10,
	}
	client := &http.Client{
		Timeout:   time.Second * 10,
		Transport: transport,
	}

	for {
		s := make(map[string][]string)
		for k, v := range Dconf.Constants {
			for r := range v {
				_, ee := client.Get(v[r])
				if ee != nil {
					fmt.Println(ee)

				} else {
					s[k] = append(s[k], v[r])
				}

			}
		}
		eq := reflect.DeepEqual(Dconf.Upstreams, s)
		switch eq {
		case false:
			Dconf.Lock()
			Dconf.Upstreams = s
			Dconf.Unlock()
		}

		d := time.Duration(healtchecks)
		time.Sleep(d * time.Second / 10)
	}

}

// curl -XPOST -u 'test:Te$ting' --data-binary @/tmp/balod.json 127.0.0.1:4141/config?cfg=new
// curl -XPOST -u 'test:Te$ting' --data-binary @/tmp/valod.json 127.0.0.1:4141/config?cfg=append
