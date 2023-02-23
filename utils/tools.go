package utils

import (
	"configs"
	"crypto/sha256"
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
	case configs.To.JWTAuth:
		jwthdr := strings.Split(r.Header.Get("Authorization"), " ")
		if CheckJWTtoken(jwthdr[len(jwthdr)-1]) {
			return true
		} else {
			http.Error(w, http.StatusText(unauth), unauth)
			log.Print("Invalid JWT authorization:")

			return false
		}
	case configs.To.BasicAuth:

		username, password, ok := r.BasicAuth()
		if ok {
			if username+":"+password == configs.To.BasicCreds {
				return true
			} else {
				log.Print("Invalid authorization:")
				http.Error(w, http.StatusText(unauth), unauth)
				return false
			}

		} else {
			log.Print("Invalid authorization:")
			http.Error(w, http.StatusText(unauth), unauth)
			return false
		}
	default:
		return false
	}
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

func CheckJWTtoken(tok string) bool {
	//tok := string(to)
	//hmacSampleSecret := []byte("Super$ecter123765@")
	hmacSampleSecret := configs.To.JWTSecret
	_, errr := jwt.Parse(tok, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return hmacSampleSecret, nil
	})
	if errr != nil {
		fmt.Println(errr)
		return false

	} else {
		//fmt.Println(token.Valid, token.Header, token.Method, token.Claims)
		return true
	}

}

// eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJoYXNoIjoiNmFjNDExOTNjYmE1ZjQ3MTZmYTQxOTAxNmZiZDJmYWQzZGJhY2M4MGU0ZDI5YWQ3ZTVlYjZkZjc1OTdiNTFjYiJ9.T1pvde1pF5hj1q9-xfmPyCtJj5qhxBOey4AXGSjKzS8
