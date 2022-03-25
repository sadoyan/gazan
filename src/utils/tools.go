package utils

import (
	"configs"
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

func CheckAuth(w http.ResponseWriter, r *http.Request) bool {

	const unauth = http.StatusUnauthorized

	switch {
	case configs.To.ApiKeyAuth:
		if r.Header.Get("X-API-KEY") == configs.To.ApiKey {
			return true
		} else {
			log.Print("Invalid API-KEY authorization:")
			http.Error(w, http.StatusText(unauth), unauth)
			return false
		}
	case configs.To.BasicAuth:
		authorized := configs.To.BasicCreds
		basicAuth := r.Header.Get("Authorization")

		if !strings.HasPrefix(basicAuth, "Basic ") {
			log.Print("Invalid authorization:", basicAuth)
			http.Error(w, http.StatusText(unauth), unauth)
			return false
		}

		up, err := base64.StdEncoding.DecodeString(basicAuth[6:])
		if err != nil {
			log.Print("authorization decode error:", err)
			http.Error(w, http.StatusText(unauth), unauth)
			return false
		}
		if string(up) != authorized {
			http.Error(w, http.StatusText(unauth), unauth)
			log.Print("authorization decode error:", err)

			return false
		}
		return true
	default:
		return false
	}

	//if configs.To.ApiKeyAuth {
	//	if r.Header.Get("X-API-KEY") == configs.To.ApiKey {
	//		return true
	//	} else {
	//		log.Print("Invalid authorization-AAAAAAA:", basicAuth)
	//		http.Error(w, http.StatusText(unauth), unauth)
	//		return false
	//	}
	//}
	//if configs.To.BasicAuth {
	//	if !strings.HasPrefix(basicAuth, "Basic ") {
	//		log.Print("Invalid authorization:", basicAuth)
	//		http.Error(w, http.StatusText(unauth), unauth)
	//		return false
	//	}
	//
	//	up, err := base64.StdEncoding.DecodeString(basicAuth[6:])
	//	if err != nil {
	//		log.Print("authorization decode error:", err)
	//		http.Error(w, http.StatusText(unauth), unauth)
	//		return false
	//	}
	//	if string(up) != authorized {
	//		http.Error(w, http.StatusText(unauth), unauth)
	//		log.Print("authorization decode error:", err)
	//
	//		return false
	//	}
	//	return true
	//}
	//return false

}

func Contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}

func CheckJWTtoken() {
	fmt.Println("")
	fmt.Println("- - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -")
	fmt.Println("Parse-Hmac")
	var hmacSampleSecret []byte
	//hmacSampleSecret := []byte("my_secret_key")

	tokenString := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJmb28iOiJiYXIiLCJmb29vb29vb29vb29vb29vb28iOjE0NDQ0Nzg0MDB9.AS9-OaZWBdI4DR_j_a0qcCP1xxTLFH52WudB2unZA14"
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		fmt.Println(hmacSampleSecret)
		return hmacSampleSecret, nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Println(claims["foooooooooooooooo"], claims["baaaaaaaaaaaaar"])
		fmt.Println(claims)

	} else {
		fmt.Println(err)
	}

	fmt.Println(*token)
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
	fmt.Println("")
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
