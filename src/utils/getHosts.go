package utils

import (
	"log"
	"net"
	"net/http"
	"time"
)

func GetHostsByHTTP() {
	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 5 * time.Second,
	}

	for {
		client := &http.Client{
			Timeout:   time.Second * 15,
			Transport: transport,
		}
		for k, v := range Serob {
			var newlist []string
			for t := range v {
				resp, ee := client.Get(v[t])
				if ee != nil {
					log.Println(ee)
				} else {
					if resp.StatusCode >= 100 && resp.StatusCode < 500 {
						resp.Body.Close()
						newlist = append(newlist, v[t])
					}
				}

			}
			switch testEq(newlist, Dconf.Upstreams[k]) {
			case false:
				Dconf.Lock()
				log.Println("Upstreams list is changes to", newlist)
				Dconf.Upstreams[k] = newlist
				Dconf.Unlock()
			}
		}
		//valod()
		time.Sleep(5 * time.Second)
	}

}

//func GetHostsByDNS() {
//	fmt.Println("")
//	iprecords, _ := net.LookupIP("dwarf.netangels.loc")
//	for _, ip := range iprecords {
//		fmt.Println(ip)
//	}
//
//	ptr, _ := net.LookupAddr("192.168.10.1")
//	for _, ptrvalue := range ptr {
//		fmt.Println(ptrvalue)
//	}
//
//	nameserver, _ := net.LookupNS("netangels.net")
//	for _, ns := range nameserver {
//		fmt.Println(ns)
//	}
//	cname, _ := net.LookupCNAME("graph.netangels.net")
//	fmt.Println(cname)
//	fmt.Println("")
//
//	txtrecords, _ := net.LookupTXT("netangels.net")
//
//	for _, txt := range txtrecords {
//		fmt.Println(txt)
//	}
//
//	//cname, srvs, err := net.LookupSRV("xmpp-server", "tcp", "golang.org")
//	//if err != nil {
//	//	log.Println(err)
//	//}
//	//
//	//fmt.Printf("\ncname: %s \n\n", cname)
//	//
//	//for _, srv := range srvs {
//	//	fmt.Printf("%v:%v:%d:%d\n", srv.Target, srv.Port, srv.Priority, srv.Weight)
//	//}
//
//}
