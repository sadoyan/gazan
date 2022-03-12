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

var pause = false

func postData(data string) {
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
	body := bytes.NewBufferString(data)
	req, _ := http.NewRequest(http.MethodPost, utils.RetRandom(), body)
	if to.clientAuth {
		req.SetBasicAuth(to.clientUser, to.clientPass)
	}

	req.Header.Add("Content-Length", strconv.Itoa(len(data)))
	resp, err := client.Do(req)
	if err != nil {
		to.queue <- data
		pause = true
		chocho <- true
		log.Println("Dead upstream:", err)
		time.Sleep(2 * time.Second)
		return
	}

	defer resp.Body.Close()
	if !(resp.StatusCode >= 200 && resp.StatusCode < 300) {
		to.queue <- data
		pause = true
		chocho <- true
		log.Println("Dead upstream:", err)
		time.Sleep(2 * time.Second)
		return
	} else {
		pause = false
		chocho <- false
	}
}
