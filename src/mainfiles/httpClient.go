package mainfiles

import (
	"bytes"
	"configs"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strings"
	"time"
	"utils"
)

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

var protohost string
var key string
var qstring string

func ProcessData(r *http.Request) (int, []uint8, error) {
	data, err := ioutil.ReadAll(r.Body)

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err == nil {
		r.Header.Add("X-FORWARDED-FOR", ip)
	}

	// Think baout this !
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(r.Body)

	if err != nil {
		log.Println(err)
		return 500, []uint8("500 Internal server error\n"), err
	}

	switch r.TLS {
	case nil:
		protohost = "http://" + r.Host
	default:
		protohost = "https://" + r.Host
	}

	switch r.URL.Path {
	case "/":
		qstring = "/"
		key = protohost
	default:
		qstring = strings.Split(r.URL.Path, "/")[1]
		key = protohost + "/" + qstring
	}

	veq, e := utils.RetRandomMap(key)

	if e == nil {
		req, histeric := http.NewRequest(r.Method, veq, bytes.NewReader(data))
		if histeric != nil {
			log.Println("Error connecting To Upstream:", veq)
			return 500, []uint8("500 Internal server error\n"), histeric
		}
		if configs.To.ClientAuth {
			req.SetBasicAuth(configs.To.ClientUser, configs.To.ClientPass)
		}
		//req.Header.Add("Content-Length", strconv.Itoa(len(data)))
		req.Header = r.Header
		resp, err := client.Do(req)
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
	} else {
		return 503, []uint8("503 Service Unavailable"), e
	}
}
