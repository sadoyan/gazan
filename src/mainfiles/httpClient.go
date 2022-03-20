package mainfiles

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
	"utils"
)

//func httpClient() *http.Client {
//	client := &http.Client{
//		Transport: &http.Transport{
//			MaxIdleConns:          100,
//			IdleConnTimeout:       90 * time.Second,
//			TLSHandshakeTimeout:   10 * time.Second,
//			ExpectContinueTimeout: 1 * time.Second,
//			MaxIdleConnsPerHost:   10,
//		},
//		Timeout: 10 * time.Second,
//	}
//	return client
//}

func postData(data map[string][]byte, method string) (int, []uint8, error) {
	//transport := &http.Transport{
	//	DialContext: (&net.Dialer{
	//		Timeout:   30 * time.Second,
	//		KeepAlive: 30 * time.Second,
	//	}).DialContext,
	//	MaxIdleConns:          100,
	//	IdleConnTimeout:       90 * time.Second,
	//	TLSHandshakeTimeout:   10 * time.Second,
	//	ExpectContinueTimeout: 1 * time.Second,
	//	MaxIdleConnsPerHost:   10,
	//}
	//client := &http.Client{
	//	Timeout:   time.Second * 10,
	//	Transport: transport,
	//}
	//client := httpClient()
	//resp, err := client.Do(req)

	//clientTrace := &httptrace.ClientTrace{
	//	GotConn: func(info httptrace.GotConnInfo) { log.Printf("conn was reused: %t", info.Reused) },
	//}
	//traceCtx := httptrace.WithClientTrace(context.Background(), clientTrace)
	//req, histeric := http.NewRequestWithContext(traceCtx, http.MethodPost, veq, bytes.NewReader(v))

	for k, v := range data {
		veq, e := utils.RetRandomMap(k)
		if e == nil {
			req, histeric := http.NewRequest(method, veq, bytes.NewReader(v))
			//req, histeric := http.NewRequestWithContext(traceCtx, http.MethodPost, veq, bytes.NewReader(v))

			if histeric != nil {
				log.Println("Error connecting to upstream:", veq)
				break
			}
			if to.clientAuth {
				req.SetBasicAuth(to.clientUser, to.clientPass)
			}
			req.Header.Add("Content-Length", strconv.Itoa(len(v)))
			resp, err := http.DefaultClient.Do(req)
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
