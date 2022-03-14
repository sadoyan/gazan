package mainfiles

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"runtime"
	"strings"
	"time"
	"utils"
)

var fullurl string

func dynHandler(w http.ResponseWriter, r *http.Request) {
	const unauth = http.StatusUnauthorized
	if to.serverAuth {
		auth := r.Header.Get("Authorization")
		if !strings.HasPrefix(auth, "Basic ") {
			log.Print("Invalid authorization:", auth)
			http.Error(w, http.StatusText(unauth), unauth)
			return
		}
		up, err := base64.StdEncoding.DecodeString(auth[6:])
		if err != nil {
			log.Print("authorization decode error:", err)
			http.Error(w, http.StatusText(unauth), unauth)
			return
		}
		if string(up) != authorized["server"] {
			http.Error(w, http.StatusText(unauth), unauth)
			return
		}
	}
	switch r.Method {
	case "POST", "PUT", "GET":
		reqBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Fatal(err)
		}
		go func(out chan<- map[string][]byte) {
			switch r.TLS {
			case nil:
				fullurl = "http://" + r.Host + r.URL.Path
			default:
				fullurl = "https://" + r.Host + r.URL.Path
			}
			m := make(map[string][]byte)
			m[fullurl] = reqBody
			out <- m
			log.Println(r.Proto, r.UserAgent(), r.RemoteAddr, r.Method, fullurl)

		}(to.queue)
	default:
		_, _ = fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}
}

func mxhandl(w http.ResponseWriter, _ *http.Request) {
	mz := printStats()
	_, _ = fmt.Fprintln(w, mz)
}
func dynconfig(w http.ResponseWriter, r *http.Request) {
	utils.ApiConfig(r)
}

func playmux0() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", dynHandler)

	s1 := http.Server{
		Addr:         to.httpAddress,
		Handler:      mux,
		ReadTimeout:  100 * time.Second,
		WriteTimeout: 100 * time.Second,
	}
	_ = s1.ListenAndServe()

}
func playmux1() {
	mux1 := http.NewServeMux()
	mux1.HandleFunc("/", mxhandl)

	s2 := http.Server{
		Addr:         "127.0.0.1:9191",
		Handler:      mux1,
		ReadTimeout:  100 * time.Second,
		WriteTimeout: 100 * time.Second,
	}
	_ = s2.ListenAndServe()

}
func playmux2() {
	mux2 := http.NewServeMux()
	mux2.HandleFunc("/config", dynconfig)

	s3 := http.Server{
		Addr:         "127.0.0.1:4141",
		Handler:      mux2,
		ReadTimeout:  100 * time.Second,
		WriteTimeout: 100 * time.Second,
	}
	_ = s3.ListenAndServe()

}

var chocho = make(chan bool)

func RunServer() {
	setVarsik()
	http.HandleFunc("/", dynHandler)
	fmt.Println("starting server at: " + to.httpAddress)

	if to.monenabled {
		go playmux1()
	}
	// ---------------------------------------------- //
	go playmux2()
	// ---------------------------------------------- //
	log.Print("Started Proxy ")
	for j := 0; j < to.dispatchersCount; j++ {
		go func() {
			for {
				s := <-to.queue
				postData(s)
			}
		}()
	}

	go func(in chan bool) {
		for {
			_ = <-in
		}
	}(chocho)

	runtime.Gosched()
	playmux0()

	//forever := make(chan bool)
	//<-forever
}
