package utils

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

func testEq(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func RetRandomMap(key string) (string, error) {

	if len(Dconf.Upstreams[key]) >= 1 {
		randomIndex := rand.Intn(len(Dconf.Upstreams[key]))
		pick := Dconf.Upstreams[key][randomIndex]
		return pick, nil

	} else {
		return "", errors.New("upstream not found, or upstreams list is empty")
	}

}

func RetRandom(input []string) string {
	randomIndex := rand.Intn(len(input))
	pick := input[randomIndex]
	return pick
}

func CheckAuth(w http.ResponseWriter, r *http.Request, authorized string) {
	const unauth = http.StatusUnauthorized
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
	if string(up) != authorized {
		http.Error(w, http.StatusText(unauth), unauth)
		return
	}
}

func CheckJWTtoken() {
	fmt.Println("- - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -")
	fmt.Println("Parse-Hmac")
	var hmacSampleSecret []byte
	//hmacSampleSecret := []byte("my_secret_key")

	tokenString := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJmb29vb29vb29vb29vb29vb28iOjE0NDQ0Nzg0MDB9.CZ0n1l35q1BUgWxHU8M7kk7u0ejh14X_lFSDoDEgSsU"
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		fmt.Println(hmacSampleSecret)
		return hmacSampleSecret, nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Println(claims["fooooooooooooooooo"], claims["baaaaaaaaaaaaar"])
	} else {
		fmt.Println(err)
	}
	fmt.Println("New-Hmac")
	// -------------------------------------------------------------------
	var hmacSampleSecret2 []byte
	token2 := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"foo":                "bar",
		"fooooooooooooooooo": time.Date(2015, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
	})
	tokenString2, err2 := token2.SignedString(hmacSampleSecret2)
	fmt.Println(tokenString2, err2)

	// -------------------------------------------------------------------

	fmt.Println("- - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -")
	//https://pkg.go.dev/github.com/golang-jwt/jwt#example-New-Hmac
}

//var seeds = []string{"https://netangels.net/utils/testurl"}
//var seed = "https://netangels.net/utils/testurl"
//var hostlist []string
//type hl struct {
//	Hostlist []string
//	sync.RWMutex
//}
//var hoho = hl{
//	Hostlist: nil,
//	RWMutex:  sync.RWMutex{},
//}
//func GetSeed() {
//	for {
//		thishosts := GetHostsByHTTP(seed)
//		if !testEq(thishosts, hostlist) {
//			hostlist = thishosts
//			hoho.RWMutex.Lock()
//			hoho.Hostlist = thishosts
//			hoho.RWMutex.Unlock()
//			fmt.Println("New hosts list is", hoho.Hostlist)
//		}
//		time.Sleep(10 * time.Second)
//	}
//
//}
