package mainfiles

import (
	"bytes"
	"log"
	"net"
	"net/http"
	"strconv"
	"time"
	"utils"
)

func postData(data map[string][]byte) {
	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
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
	for k, v := range data {
		veq, e := utils.RetRandomMap(k)
		if e == nil {
			req, histeric := http.NewRequest(http.MethodPost, veq, bytes.NewReader(v))
			if histeric != nil {
				log.Println("Error connecting to upstream:", veq)
				break
			}
			if to.clientAuth {
				req.SetBasicAuth(to.clientUser, to.clientPass)
			}
			req.Header.Add("Content-Length", strconv.Itoa(len(v)))
			resp, err := client.Do(req)

			if err != nil {
				to.queue <- data

				chocho <- true
				log.Println("Dead upstream:", err)
				time.Sleep(2 * time.Second)
				return
			}
			if !(resp.StatusCode >= 200 && resp.StatusCode < 300) {
				to.queue <- data

				chocho <- true
				log.Println("Dead upstream:", err)
				time.Sleep(2 * time.Second)
				return
			} else {

				chocho <- false
			}
			_ = resp.Body.Close()
		} else {
			log.Println(e)
		}

	}
	// ------------------------------------ //

}
