package utils

import (
	"fmt"
	"net"
)

//func GetHostsByHTTP(seed string) []string {
//	transport := &http.Transport{
//		DialContext: (&net.Dialer{
//			Timeout:   30 * time.Second,
//			KeepAlive: 30 * time.Second,
//		}).DialContext,
//		MaxIdleConns:          100,
//		IdleConnTimeout:       90 * time.Second,
//		TLSHandshakeTimeout:   10 * time.Second,
//		ExpectContinueTimeout: 5 * time.Second,
//	}
//	client := &http.Client{
//		Timeout:   time.Second * 15,
//		Transport: transport,
//	}
//	resp, ee := client.Get(seed)
//	if ee != nil {
//		fmt.Println(ee)
//		return hostlist
//	}
//	defer resp.Body.Close()
//	bodyBytes, _ := ioutil.ReadAll(resp.Body)
//	hosts := strings.Split(string(bodyBytes), "\n")
//
//	var thishosts []string
//
//	for host := range hosts {
//		if hosts[host] != "" {
//			thishosts = append(thishosts, hosts[host])
//		}
//	}
//
//	return thishosts
//
//}

func GetHostsByDNS() {
	fmt.Println("")
	iprecords, _ := net.LookupIP("dwarf.netangels.loc")
	for _, ip := range iprecords {
		fmt.Println(ip)
	}

	ptr, _ := net.LookupAddr("192.168.10.1")
	for _, ptrvalue := range ptr {
		fmt.Println(ptrvalue)
	}

	nameserver, _ := net.LookupNS("netangels.net")
	for _, ns := range nameserver {
		fmt.Println(ns)
	}
	cname, _ := net.LookupCNAME("graph.netangels.net")
	fmt.Println(cname)
	fmt.Println("")

	txtrecords, _ := net.LookupTXT("netangels.net")

	for _, txt := range txtrecords {
		fmt.Println(txt)
	}

	//cname, srvs, err := net.LookupSRV("xmpp-server", "tcp", "golang.org")
	//if err != nil {
	//	log.Println(err)
	//}
	//
	//fmt.Printf("\ncname: %s \n\n", cname)
	//
	//for _, srv := range srvs {
	//	fmt.Printf("%v:%v:%d:%d\n", srv.Target, srv.Port, srv.Priority, srv.Weight)
	//}

}
