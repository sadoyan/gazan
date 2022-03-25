package mainfiles

import (
	"bytes"
	"configs"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strconv"
	"time"
	"utils"
)

//var transport = &http.Transport{
//	DialContext: (&net.Dialer{
//		Timeout:   30 * time.Second,
//		KeepAlive: 30 * time.Second,
//	}).DialContext,
//	MaxIdleConns:          100,
//	MaxConnsPerHost:       10,
//	MaxIdleConnsPerHost:   10,
//	IdleConnTimeout:       90 * time.Second,
//	TLSHandshakeTimeout:   10 * time.Second,
//	ExpectContinueTimeout: 1 * time.Second,
//}
//var client = &http.Client{
//	Timeout:   time.Second * 10,
//	Transport: transport,
//}

var transport = &http.Transport{
	DialContext: (&net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}).DialContext,
	MaxIdleConns:          configs.To.Clienmaxidle,
	MaxConnsPerHost:       configs.To.Clienmaxperhost,
	MaxIdleConnsPerHost:   configs.To.Clienmaxidleperhost,
	IdleConnTimeout:       configs.To.Clienidletimeout * time.Second,
	TLSHandshakeTimeout:   10 * time.Second,
	ExpectContinueTimeout: 1 * time.Second,
}
var client = &http.Client{
	Timeout:   time.Second * configs.To.Clientimeout,
	Transport: transport,
}

func PostData(data map[string][]byte, method string) (int, []uint8, error) {
	for k, v := range data {
		veq, e := utils.RetRandomMap(k)
		if e == nil {
			req, histeric := http.NewRequest(method, veq, bytes.NewReader(v))
			//req, histeric := http.NewRequestWithContext(traceCtx, http.MethodPost, veq, bytes.NewReader(v))

			if histeric != nil {
				log.Println("Error connecting To Upstream:", veq)
				break
			}
			if configs.To.ClientAuth {
				req.SetBasicAuth(configs.To.ClientUser, configs.To.ClientPass)
			}
			req.Header.Add("Content-Length", strconv.Itoa(len(v)))
			//resp, err := http.DefaultClient.Do(req)
			resp, err := client.Do(req)
			// - - - - - - - - - - - - - - - - - - - - - - -
			if err != nil {
				log.Println("Dead upstream:", err)
				time.Sleep(2 * time.Second)
				return 500, nil, err
			}

			buf, buerr := ioutil.ReadAll(resp.Body)
			if buerr != nil {
				log.Println("clientient read body error", buerr)
			}

			if !(resp.StatusCode >= 200 && resp.StatusCode < 300) {
				log.Println("Dead upstream:", err)
				time.Sleep(2 * time.Second)
				fmt.Println("")
				_ = resp.Body.Close()
				return resp.StatusCode, nil, err
			} else {
				_ = resp.Body.Close()
				return resp.StatusCode, buf, nil
			}
			//_ = resp.Body.Close()
		} else {
			return 503, []uint8("503 Service Unavailable"), e
			//log.Println(e)
		}

	}
	return 500, nil, nil
}
