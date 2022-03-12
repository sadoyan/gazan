package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

//var dynconfigs []map[string]interface{}
var dynconfigs map[string][]string

type ddconf struct {
	coco map[string][]string
}

var dconf *ddconf

func logprint(msg string, err error) {
	if err != nil {
		log.Println(msg, err)
	}
}

func ApiConfig(r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	er := decoder.Decode(&dynconfigs)
	logprint("Json Decode error:", er)

	for k, v := range dynconfigs {
		fmt.Println("Key:", k)
		for vv := range v {
			fmt.Println("Value:", v[vv])
		}
		fmt.Println(" ")
	}

	//fmt.Println("- - - - - - - - - - - - - - - - - ")
	//decot := json.Marshal()
	//e := decoder.Decode(dconf.coco)
	//logprint("Json Decode error:", e)

	//for k, v := range dconf.coco {
	//	fmt.Println("Key:", k)
	//	for vv := range v {
	//		fmt.Println("Value:", v[vv])
	//	}
	//	fmt.Println(" ")
	//}

}
