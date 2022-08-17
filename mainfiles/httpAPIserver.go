package mainfiles

import (
	"configs"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"runtime"
	"time"
	"utils"
)

func dynHandler(w http.ResponseWriter, r *http.Request) {
	c.inc()

	if configs.To.ServerAuth {
		if !utils.CheckAuth(w, r) {
			return
		}
	}
	switch r.Method {
	case "POST", "GET":

		status, body, headers, err := ProcessData(r, w)

		if err != nil {

			w.WriteHeader(status)
			_, be := w.Write([]uint8("500 Internal server error\n"))
			if be != nil {
				log.Println(be)
			}
		}

		for k, v := range headers {
			for hlen := range v {
				w.Header().Add(k, v[hlen])
			}

		}
		//w.WriteHeader(status)
		_, ee := w.Write(body)
		if ee != nil {
			log.Println("HTTP basic error:", ee)
			_, _ = w.Write([]uint8("500 Internal server error\n"))
		}
		if configs.To.Accesslog {
			log.Println(status, r.Proto, r.RemoteAddr, r.Method, r.Host+r.RequestURI)
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
	_, ee := w.Write(utils.ApiConfig(r))
	if ee != nil {
		log.Println("Error in API config", ee)
	}
	_, _ = w.Write([]byte("\n"))
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
	mux.HandleFunc("/loginski", jwtLogin)
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
		Addr:         configs.To.Monurl,
		Handler:      mux1,
		ReadTimeout:  100 * time.Second,
		WriteTimeout: 100 * time.Second,
	}
	fmt.Println("Starting monitoring at:", configs.To.Monurl)
	_ = s2.ListenAndServe()
}
func playmux2() {
	mux2 := http.NewServeMux()
	mux2.HandleFunc("/config", dynconfig)

	s3 := http.Server{
		Addr:         configs.To.Confurl,
		Handler:      mux2,
		ReadTimeout:  100 * time.Second,
		WriteTimeout: 100 * time.Second,
	}
	fmt.Println("Starting server at:", configs.To.HttpAddress)
	_ = s3.ListenAndServe()

}

func serveTLS() {
	muxTLS := http.NewServeMux()
	muxTLS.HandleFunc("/", dynHandler)
	muxTLS.HandleFunc("/loginski", jwtLogin)
	sTLS := http.Server{
		Addr:         configs.To.TLSAddress,
		Handler:      muxTLS,
		ReadTimeout:  100 * time.Second,
		WriteTimeout: 100 * time.Second,
	}

	r, e1 := ioutil.ReadFile(configs.To.TLSCertFIle)
	block, _ := pem.Decode(r)

	if e1 != nil {
		log.Fatal("Error loading certificate file:", e1)
	}

	_, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		log.Fatal("Invalid Certificate file:", err)
	}

	v, e3 := ioutil.ReadFile(configs.To.TLSPrivKey)
	if e3 != nil {
		log.Fatal("Error loading private key file:", e3)
	}

	b, _ := pem.Decode(v)

	_, er := x509.ParsePKCS1PrivateKey(b.Bytes)
	if er != nil {
		log.Fatal("Invalid private key file", er)
	}
	fmt.Println("Starting TLS server at:", configs.To.TLSAddress)
	_ = sTLS.ListenAndServeTLS(configs.To.TLSCertFIle, configs.To.TLSPrivKey)
}

func RunServer() {
	configs.SetVarsik()
	http.HandleFunc("/", dynHandler)

	if configs.To.Monenabled {
		go playmux1()
	}
	go playmux2()
	runtime.Gosched()
	go playmux0()

	if configs.To.TLSEnabled {
		go serveTLS()
	}
	time.Sleep(time.Second)
	utils.LoadUpstreamsFronFIle(configs.To.UpstreamsFile)
	forever := make(chan bool)
	<-forever

}
