package utils

import (
	"configs"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
)

type auth interface {
	auth() bool
}

type basic struct {
	User string
	Pass string
	Auth bool
}
type token struct {
	token string
}
type api struct {
	Key string
}

type credential struct {
	Keys  map[string]bool
	User  map[string]string
	Token []byte
	sync.RWMutex
}

var Credential = &credential{
	Keys:    map[string]bool{},
	User:    map[string]string{},
	Token:   []byte{},
	RWMutex: sync.RWMutex{},
}

func PopulateUSers() {
	Credential.Lock()
	Credential.Keys[os.Getenv("GAZANKEY")] = true
	up := md5.Sum([]byte(os.Getenv("BASICPASS")))
	Credential.User[os.Getenv("BASICUSER")] = hex.EncodeToString(up[:])
	Credential.Token = []byte(os.Getenv("JWTSECRET"))
	Credential.Unlock()
}

func (ba *basic) auth() bool {
	md5HashInBytes := md5.Sum([]byte(ba.Pass))
	md5HashInString := hex.EncodeToString(md5HashInBytes[:])
	pass, ok := Credential.User[ba.User]

	if !ok {
		return false
	}
	return md5HashInString == pass
}
func (ap *api) auth() bool {
	if _, ok := Credential.Keys[ap.Key]; ok {
		return true
	} else {
		return false
	}
}
func (tk *token) auth() bool {
	tok := tk.token
	hmacSampleSecret := Credential.Token
	_, errr := jwt.Parse(tok, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return hmacSampleSecret, nil
	})
	if errr != nil {
		return false
	} else {
		return true
	}

}

func CheckAuth(w http.ResponseWriter, r *http.Request) bool {
	switch {
	case configs.To.ApiKeyAuth:
		a := api{Key: r.Header.Get("X-API-KEY")}
		if !a.auth() {
			log.Print("Invalid API-KEY authorization:")
			http.Error(w, http.StatusText(403), 403)
		}
		return a.auth()
	case configs.To.BasicAuth:
		username, password, ok := r.BasicAuth()
		if !ok {
			return false
		}
		b := basic{User: username, Pass: password, Auth: ok}
		x := b.auth()
		if !x {
			log.Print("Invalid user/pass for basic auth:")
			http.Error(w, http.StatusText(403), 403)
		}
		return x
	case configs.To.JWTAuth:
		jwthdr, ok := r.URL.Query()["token"]
		if !ok {
			jwthdr = strings.Split(r.Header.Get("Authorization"), " ")
		}
		c := token{token: jwthdr[len(jwthdr)-1]}
		yes := c.auth()
		if !yes {
			log.Print("Invalid JWT Token:")
			http.Error(w, http.StatusText(403), 403)
		}
		return yes
	}
	return false
}

type jwtinput struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Expire   int    `json:"exp"`
}

func GenJWTtoken(in []byte) ([]byte, error) {
	var jwtin jwtinput
	err := json.Unmarshal(in, &jwtin)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	hmacSampleSecret := Credential.Token
	hash := fmt.Sprintf("%x", sha256.Sum256([]byte(jwtin.Username+jwtin.Password)))
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"hash": hash,
		"exp":  jwtin.Expire,
	})

	tokenString, err2 := claims.SignedString(hmacSampleSecret)
	if err != nil {
		log.Println("Error Getting JWT signed key:", err2)
		return nil, err2
	}
	fmt.Println(jwtin.Username, jwtin.Password, jwtin.Expire, tokenString)
	return []byte(tokenString), nil
	// curl -s -XPOST -d'{"username":"uzver", "password":"PazZik", "exp":1682003277}' 192.168.10.10:8080/loginski
}
