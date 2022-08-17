package mainfiles

import (
	"encoding/json"
	"runtime"
	"sync/atomic"
	"time"
)

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

type count32 int32

func (c *count32) inc() int32 {
	return atomic.AddInt32((*int32)(c), 1)
}

func (c *count32) get() int32 {
	return atomic.LoadInt32((*int32)(c))
}

var m runtime.MemStats
var c count32
var b int32 = 1
var t int32 = 1

type metrics struct {
	Alloc           uint64 `json:"alloc,int"`
	Total           uint64 `json:"total,int"`
	System          uint64 `json:"system,int"`
	Gcnum           uint32 `json:"gcnum,int"`
	Frees           uint64 `json:"frees,int"`
	HeapAlloc       uint64 `json:"heapalloc,int"`
	HeapIdle        uint64 `json:"heapidle,int"`
	HeapInuse       uint64 `json:"heapinuse,int"`
	HeapObjects     uint64 `json:"heapobjects,int"`
	HeapReleased    uint64 `json:"heapreleased,int"`
	LastGC          uint64 `json:"lastgc,int"`
	NumForcedGC     uint32 `json:"forcegc,int"`
	PauseTotalNs    uint64 `json:"pausetotal,int"`
	Goroutines      int    `json:"goroutines,int"`
	AccessCounter   int32  `json:"accesscounter,int"`
	AccessPerSecond int32  `json:"accesspersecond,int"`
}

func printStats() (s string) {
	runtime.ReadMemStats(&m)
	u := &metrics{}
	u.Alloc = m.Alloc
	u.Total = m.TotalAlloc
	u.System = m.Sys
	u.Gcnum = m.NumGC
	u.Frees = m.Frees
	u.HeapAlloc = m.HeapAlloc
	u.HeapIdle = m.HeapIdle
	u.HeapInuse = m.HeapInuse
	u.HeapObjects = m.HeapObjects
	u.HeapReleased = m.HeapReleased
	u.PauseTotalNs = m.PauseTotalNs
	u.NumForcedGC = m.NumForcedGC
	u.Goroutines = runtime.NumGoroutine()
	u.AccessCounter = c.get()
	tn := int32(time.Now().Unix())
	if u.AccessCounter-b > 0 && tn-t > 0 {
		u.AccessPerSecond = (u.AccessCounter - b) / (tn - t)
	}
	b = u.AccessCounter
	t = tn
	result, _ := json.MarshalIndent(u, "", "    ")

	// -------------------------------------------------------- //
	//dnsServer := "1.1.1.1:53"
	//
	//r := &net.Resolver{
	//	PreferGo: true,
	//	Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
	//		d := net.Dialer{
	//			Timeout: time.Millisecond * time.Duration(10000),
	//		}
	//		return d.DialContext(ctx, network, dnsServer)
	//	},
	//}
	//ip, _ := r.LookupHost(context.Background(), "app.oddeye.co")
	//tx, _ := r.LookupTXT(context.Background(), "netangels.net")
	//
	//fmt.Println(ip)
	//fmt.Println(tx)
	//
	//xx := len(dnsServer)
	//
	//cname, srvs, err := net.LookupSRV("xmpp-server", "tcp", "google.com")
	//
	//if xx != 0 {
	//	cname, srvs, err = r.LookupSRV(context.Background(), "xmpp-server", "tcp", "google.com")
	//}
	//
	//if err != nil {
	//	log.Panic(err)
	//}
	//fmt.Println(cname)
	//for _, srv := range srvs {
	//	fmt.Println(srv.Target, srv.Port, srv.Priority, srv.Weight)
	//}
	// -------------------------------------------------------- //
	return string(result)
}
