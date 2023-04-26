package mainfiles

import (
	"bufio"
	"bytes"
	"configs"
	"fmt"
	"io"
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

var key string
var qstring string
var v string

func ProcessData(r *http.Request, w http.ResponseWriter) (int, []uint8, http.Header, error) {
	data, err := io.ReadAll(r.Body)
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err == nil {
		r.Header.Add("X-FORWARDED-FOR", ip)
	}

	defer r.Body.Close()

	if err != nil {
		log.Println(err)
		return 500, []uint8("500 Internal server error\n"), nil, err
	}

	switch utils.Dconf.Windcards[r.Host] {
	case true:
		qstring = strings.Split(r.URL.Path, "/")[1]
		z := r.Host + "/" + qstring
		if _, ok := utils.Dconf.Upstreams[z]; ok {
			key = r.Host + "/" + qstring
		} else {
			key = r.Host
		}

		veq, e := utils.RetRandomMap(key)
		v = veq + r.URL.String()
		if e != nil {
			return 503, []uint8("503 Service Unavailable\n"), nil, e
		}
	default:
		key = r.Host
		q := strings.Split(r.URL.Path, "/")
		var veq string
		for x := range q {
			if x != 0 {
				key = key + "/" + q[x]
				veq, err = utils.RetRandomMap(key)
				v = veq + r.URL.String()
				if err == nil {
					break
				}
			}
		}
		if len(veq) == 0 {
			return 503, []uint8("503 Service Unavailable\n"), nil, err
		}
	}

	req, histeric := http.NewRequest(r.Method, v, bytes.NewReader(data))
	if histeric != nil {
		log.Println("Error connecting To Upstream:")
		return 500, []uint8("500 Internal server error\n"), nil, histeric
	}
	if configs.To.ClientAuth {
		req.SetBasicAuth(configs.To.ClientUser, configs.To.ClientPass)
	}

	req.Header = r.Header
	req.Host = r.Host
	resp, err := client.Do(req)

	if resp.Request.Response != nil {
		resp.StatusCode = resp.Request.Response.StatusCode
		resp.Header.Add("Location", resp.Request.URL.Path)
	}
	if err != nil {
		log.Println("Dead upstream:", err)
		time.Sleep(2 * time.Second)
		return 500, nil, nil, err
	}
	if resp.ContentLength > 1048576 {
		scanner := bufio.NewScanner(resp.Body)
		scanner.Split(bufio.ScanBytes)
		for k, v := range resp.Header {
			for x := range v {
				w.Header().Add(k, v[x])
			}
		}
		for scanner.Scan() {
			_, ee := w.Write(scanner.Bytes())
			if ee != nil {
				log.Println("Error downloading big file", ee)
				return 503, []uint8("500 Service Unavailable\n"), nil, ee
			}
		}
		return resp.StatusCode, nil, nil, nil
	}
	buf, buerr := io.ReadAll(resp.Body)
	if buerr != nil {
		log.Println("clientient read body error", buerr)
	}
	if !(resp.StatusCode >= 100 && resp.StatusCode <= 500) {
		log.Println("Dead upstream:", err)
		time.Sleep(2 * time.Second)
		fmt.Println("")
		_ = resp.Body.Close()
		return resp.StatusCode, nil, nil, err
	} else {
		_ = resp.Body.Close()
		return resp.StatusCode, buf, resp.Header, nil
	}
}
