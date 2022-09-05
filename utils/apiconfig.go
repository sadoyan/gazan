package utils

import (
	"configs"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"
)

func wcconfig(k, s string) {
	if strings.HasSuffix(k, "/*") {
		nk := strings.Replace(k, "/*", "", -1)
		Dconf.Lock()
		Dconf.Windcards[nk] = true
		Dconf.Constants[nk] = Dconf.Constants[k]
		Dconf.Upstreams[nk] = Dconf.Upstreams[k]
		delete(Dconf.Constants, k)
		delete(Dconf.Upstreams, k)
		Dconf.Unlock()
	} else if !strings.Contains(k, s) {
		Dconf.Lock()
		Dconf.Windcards[k] = true
		Dconf.Unlock()
	} else {

	}
}

func ApiConfig(r *http.Request) []byte {
	urlparam := r.URL.Query().Get("cfg")
	decoder := json.NewDecoder(r.Body)
	switch urlparam {
	case "get", "dump":
		result, _ := json.MarshalIndent(Dconf, "", "    ")
		return result
	case "append", "add":
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
			wcconfig(k, "/")
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
			wcconfig(k, "/")
		}
	default:
		log.Println("Unknown parameter ")
	}
	result, _ := json.MarshalIndent(Dconf, "", "    ")
	return result
}

func LoadUpstreamsFronFIle(up string) {
	data, err := ioutil.ReadFile(up)
	if err != nil {
		log.Println("Cant load default upstreams file:", err)
		log.Println("Startingwithout upstreams")
	} else {
		er := json.Unmarshal(data, &Dconf.Upstreams)
		for k, v := range Dconf.Upstreams {

			Dconf.Constants[k] = v
			for vv := range v {
				fmt.Println("Registering URL", k, "To Upstream:", v[vv])
			}
			fmt.Println(" ")
			wcconfig(k, "/")
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
					log.Println(ee)
				} else {
					s[k] = append(s[k], v[r])
				}

			}
		}
		switch reflect.DeepEqual(Dconf.Upstreams, s) {
		case false:
			Dconf.Lock()
			Dconf.Upstreams = s
			Dconf.Unlock()
		}

		d := time.Duration(healtchecks)
		time.Sleep(d * time.Second / 10)
	}

}

// -------------------------------------------------------- //

var r = &net.Resolver{
	PreferGo: true,
	Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
		d := net.Dialer{
			Timeout: time.Millisecond * time.Duration(1000),
		}
		return d.DialContext(ctx, network, RetRandom(configs.To.DnsServer))
	},
}

func choosedns(n string, p string, d string) ([]*net.SRV, error) {
	switch configs.To.Dns {
	case true:
		_, srvs, err := r.LookupSRV(context.Background(), n, p, d)
		return srvs, err
	default:
		_, srvs, err := net.LookupSRV(n, p, d)
		return srvs, err
	}
}

func GetHostsbyDNS() {
	for !configs.To.ImRunning {
		time.Sleep(500 * time.Millisecond)
	}
	switch configs.To.Dns {
	case true:
		for {
			tempUps := make(map[string][]string)
			for name, _ := range configs.To.DnsRecords {
				lookupstring := configs.To.DnsRecords[name]
				lookupslice := strings.Split(lookupstring, ".")
				n := lookupslice[0]
				p := lookupslice[1]
				d := strings.Join(lookupslice[2:], ".")
				srvs, err := choosedns(n, p, d)
				if err == nil {
					for _, srv := range srvs {
						ip, _ := r.LookupHost(context.Background(), srv.Target)
						target := "http://" + ip[0] + ":" + strconv.Itoa(int(srv.Port))
						tempUps[name] = append(tempUps[name], target)
					}
				} else {
					log.Println(err)
				}
			}
			for name, _ := range configs.To.DnsRecords {
				wcconfig(name, "/")
			}
			for _, mo := range tempUps {
				sort.Strings(mo)
			}
			if !reflect.DeepEqual(Dconf.Constants, tempUps) || !reflect.DeepEqual(Dconf.Upstreams, tempUps) {
				Dconf.Lock()
				Dconf.Constants = tempUps
				Dconf.Upstreams = tempUps
				Dconf.Unlock()
			}
			time.Sleep(5 * time.Second)
		}
	}
}

// -------------------------------------------------------- //

// curl -XPOST -u 'test:Te$ting' --data-binary @/tmp/balod.json 127.0.0.1:4141/config?cfg=new
// curl -XPOST -u 'test:Te$ting' --data-binary @/tmp/valod.json 127.0.0.1:4141/config?cfg=append
// curl -u 'test:Te$ting' 127.0.0.1:4141/config?cfg=get
