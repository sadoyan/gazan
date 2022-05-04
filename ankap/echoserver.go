package ankap

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

func mxhandl(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if len(data) > 0 {
		fmt.Println(string(data), err)
		fmt.Println("Headers:")
		for hk, hv := range r.Header {
			fmt.Println(hk, hv)
		}
		fmt.Println("Cookies:")
		for ck, cv := range r.Cookies() {
			fmt.Println(ck, cv)
		}
		fmt.Println(r.URL)
		fmt.Println("")
	}

}

func Play() {
	mux1 := http.NewServeMux()
	mux1.HandleFunc("/", mxhandl)

	s2 := http.Server{
		Addr:         "0.0.0.0:8000",
		Handler:      mux1,
		ReadTimeout:  100 * time.Second,
		WriteTimeout: 100 * time.Second,
	}
	_ = s2.ListenAndServe()

}
