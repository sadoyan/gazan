package utils

import (
	"configs"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"log"
	"math/rand"
	"net/http"
	"strings"
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

type jwtinput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func GenJWTtoken(in []byte) ([]byte, error) {
	var jwtin jwtinput
	err := json.Unmarshal(in, &jwtin)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	hmacSampleSecret := []byte("Super$ecter123765@")
	//token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
	//	"username": jwtin.Username,
	//	"password": jwtin.Password,
	//})
	hash := fmt.Sprintf("%x", sha256.Sum256([]byte(jwtin.Username+jwtin.Password)))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"hash": hash,
	})

	tokenString, err2 := token.SignedString(hmacSampleSecret)
	if err != nil {
		log.Println("Error Getting JWT signed key:", err2)
		return nil, err2
	}
	//log.Println("Generating new JWT token for", jwtin.Username)
	fmt.Println(jwtin.Username, jwtin.Password, tokenString)
	return []byte(tokenString), nil
	// curl -XPOST -d '{"username":"gesho", "password": "polozmukuck"}' http://127.0.0.1:8080/login
}

func CheckJWTtoken() {
	fmt.Println("")
	fmt.Println("- - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -")
	//fmt.Println("Parse-Hmac")
	//hmacSampleSecret := []byte("my_secret_key")
	//
	//tokenString := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJmb28iOiJiYXIiLCJmb29vb29vb29vb29vb29vb28iOjE0NDQ0Nzg0MDB9.AS9-OaZWBdI4DR_j_a0qcCP1xxTLFH52WudB2unZA14"
	//token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
	//	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
	//		return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
	//	}
	//	fmt.Println(hmacSampleSecret)
	//	return hmacSampleSecret, nil
	//})
	//
	//if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
	//	fmt.Println(claims["foooooooooooooooo"], claims["baaaaaaaaaaaaar"])
	//	fmt.Println(claims)
	//
	//} else {
	//	fmt.Println(err)
	//}
	//
	//fmt.Println(*token)
	//fmt.Println("New-Hmac")
	// -------------------------------------------------------------------
	//var hmacSampleSecret2 []byte
	hmacSampleSecret2 := []byte("Super$ecter123765@")
	token2 := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": "valod",
		"password": "Rembo3Rembo4",
	})
	tokenString2, err2 := token2.SignedString(hmacSampleSecret2)
	fmt.Println(hmacSampleSecret2, tokenString2, err2)

	// -------------------------------------------------------------------

	fmt.Println("- - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -")
	fmt.Println("")
}
