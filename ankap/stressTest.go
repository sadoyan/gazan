package ankap

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

var transport = &http.Transport{
	DialContext: (&net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}).DialContext,
	MaxIdleConns:          10,
	MaxConnsPerHost:       20,
	MaxIdleConnsPerHost:   20,
	IdleConnTimeout:       3 * time.Second,
	TLSHandshakeTimeout:   10 * time.Second,
	ExpectContinueTimeout: 1 * time.Second,
}
var client = &http.Client{
	Timeout:   time.Second * 2,
	Transport: transport,
}

var file = "/tmp/list.txt"
var requests = flag.Int("r", 0, "an int")
var parallel = flag.Int("p", 1, "an int")
var data = []byte(`{"fname":"4bad3356-4abd-4d18","mname":"e935983a-caeb-4e6b","sname":"eb3e2c43-9a20-457a"}`)
var method = "POST"

func retRandom(input []string) string {
	randomIndex := rand.Intn(len(input))
	pick := input[randomIndex]
	return pick
}

func runit(urls []string) {
	for i := 1; i <= *requests; i++ {
		req, h := http.NewRequest(method, retRandom(urls), bytes.NewReader(data))
		if h != nil {
			fmt.Println("Request Error:", h)
		}
		fo, _ := client.Do(req)
		_, e := ioutil.ReadAll(fo.Body)
		if e != nil {
			fmt.Println("Response Error:", e)
		}
	}
}

func main() {
	flag.Parse()
	fmt.Println(" ")

	fmt.Println("Running:", *requests, "requests in", *parallel, "threads")

	u, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Print(err)
	}
	var urls []string
	t := strings.Split(string(u), "\n")
	for x := range t {
		if t[x] != "" {
			urls = append(urls, t[x])
		}
	}
	startTime := time.Now()
	var wg sync.WaitGroup

	for i := 1; i <= *parallel; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			runit(urls)
		}()
	}
	wg.Wait()
	elapsed := time.Since(startTime)
	totalRequests := *parallel * *requests
	fmt.Println(" ")
	fmt.Println("Duratioin  :", elapsed.Seconds())
	fmt.Println("Executed   :", totalRequests)
	fmt.Println("Per Second :", int(float64(totalRequests)/elapsed.Seconds()))
	fmt.Println(" ")
}
