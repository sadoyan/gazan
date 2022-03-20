package mainfiles

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"runtime"
	"time"
	"utils"
)

func dynHandler(w http.ResponseWriter, r *http.Request) {
	var fullurl string
	if to.serverAuth {
		utils.CheckAuth(w, r, authorized["server"])
	}
	switch r.Method {
	case "POST", "GET":
		reqBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
		}
		switch r.TLS {
		case nil:
			fullurl = "http://" + r.Host + r.URL.Path
		default:
			fullurl = "https://" + r.Host + r.URL.Path
		}
		m := make(map[string][]byte)
		m[fullurl] = reqBody
		status, body, err := postData(m, r.Method)
		w.WriteHeader(status)
		_, ee := w.Write(body)
		if ee != nil {
			log.Println(ee)
		}
		log.Println(r.Proto, r.RemoteAddr, r.Method, fullurl)
	default:
		w.WriteHeader(501)
		_, ee := w.Write([]byte("Method (" + r.Method + ") Not implemented"))
		if ee != nil {
			log.Println(ee)
		}
	}
}

func mxhandl(w http.ResponseWriter, _ *http.Request) {
	mz := printStats()
	_, _ = fmt.Fprintln(w, mz)
}
func dynconfig(w http.ResponseWriter, r *http.Request) {
	if to.serverAuth {
		utils.CheckAuth(w, r, authorized["server"])
	}
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

func RunServer() {
	setVarsik()
	http.HandleFunc("/", dynHandler)
	fmt.Println("starting server at: " + to.httpAddress)

	if to.monenabled {
		go playmux1()
	}
	go playmux2()
	log.Print("Started Proxy ")
	runtime.Gosched()
	go playmux0()
	time.Sleep(time.Second)
	utils.LoadUpstreams(to.upstreamsFile)
	forever := make(chan bool)
	<-forever

}
