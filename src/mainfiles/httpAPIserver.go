package mainfiles

import (
	"configs"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"runtime"
	"strings"
	"time"
	"utils"
)

func dynHandler(w http.ResponseWriter, r *http.Request) {
	if configs.To.ServerAuth {
		if !utils.CheckAuth(w, r) {
			return
		}
	}

	// -- -- JWT Testing -- -- //
	jwtheader := r.Header.Get("Authorization")
	const unauth = http.StatusUnauthorized
	if !strings.HasPrefix(jwtheader, "Authorization ") {
		jwt := r.Header.Get("Authorization")
		if !utils.CheckJWTtoken(jwt[7:]) {
			http.Error(w, http.StatusText(unauth), unauth)
			return
		}
	}
	// -- -- JWT Testing -- -- //

	switch r.Method {
	case "POST", "GET":
		//reqBody, err := ioutil.ReadAll(r.Body)
		//if !utils.CheckJWTtoken(reqBody) {
		//	w.WriteHeader(http.StatusUnauthorized)
		//	w.Write([]byte("Fuck off\n"))
		//	return
		//}
		//if err != nil {
		//	log.Println(err)
		//}
		//switch r.TLS {
		//case nil:
		//	fullurl = "http://" + r.Host + r.RequestURI
		//	//host = "http://" + r.Host
		//default:
		//	fullurl = "https://" + r.Host + r.RequestURI
		//	//host = "https://" + r.Host
		//}
		//m := make(map[string][]byte)
		//m[fullurl] = reqBody
		//fmt.Println("Host:      ", r.Host)
		//fmt.Println("RawQuery:  ", r.URL.RawQuery)
		//fmt.Println("URL:       ", r.URL)
		//fmt.Println("URL.Path:  ", r.URL.Path)
		//fmt.Println("RequestURI:", r.RequestURI)
		//fmt.Println("URL.String:", r.URL.String())
		//fmt.Println("FullURL:   ", fullurl)
		//status, body, err := PostData(m, r.Method)
		//status, body, err := ProcessData(fullurl, reqBody, r.Method)
		status, body, err := ProcessData(r)
		if err != nil {
			w.WriteHeader(status)
			_, be := w.Write([]uint8("500 Internal server error\n"))
			if be != nil {
				log.Println(be)
			}
		}
		w.WriteHeader(status)
		_, ee := w.Write(body)
		if ee != nil {
			log.Println(ee)
		}
		if configs.To.Accesslog {
			log.Println(r.Proto, r.RemoteAddr, r.Method, r.Host+r.RequestURI)
		}
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
	if configs.To.ServerAuth {
		utils.CheckAuth(w, r)
	}
	utils.ApiConfig(r)
}
func jwtLogin(w http.ResponseWriter, r *http.Request) {
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("JWT handler read body", err)
	}
	tok, er := utils.GenJWTtoken(reqBody)
	if er != nil {
		w.WriteHeader(503)
		_, ee := w.Write([]byte("Error decoding JWT token\n"))
		if ee != nil {
			log.Println(ee)
		}
	} else {
		w.WriteHeader(200)
		_, ee := w.Write(tok)
		if ee != nil {
			log.Println(ee)
		}
		_, _ = w.Write([]byte("\n"))
	}
}

func playmux0() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", dynHandler)
	mux.HandleFunc("/login", jwtLogin)
	s1 := http.Server{
		Addr:         configs.To.HttpAddress,
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
	configs.SetVarsik()
	http.HandleFunc("/", dynHandler)
	fmt.Println("starting server at: " + configs.To.HttpAddress)

	if configs.To.Monenabled {
		go playmux1()
	}
	go playmux2()
	log.Print("Started Proxy ")
	runtime.Gosched()
	go playmux0()
	time.Sleep(time.Second)
	utils.LoadUpstreamsFronFIle(configs.To.UpstreamsFile)
	forever := make(chan bool)
	<-forever

}
