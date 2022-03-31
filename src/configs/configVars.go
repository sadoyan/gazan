package configs

import (
	"flag"
	"fmt"
	"gopkg.in/ini.v1"
	"log"
	"os"
	"strings"
	"time"
)

func strtoBool(in string) bool {
	out := false
	if strings.ToLower(in) == "on" || strings.ToLower(in) == "yes" {
		out = true
	}
	return out
}

type confVars struct {
	HttpAddress         string
	Healtchecks         int
	ServerAuth          bool
	BasicAuth           bool
	ApiKeyAuth          bool
	JWTAuth             bool
	ApiKey              string
	JWTSecret           []byte
	BasicCreds          string
	ClientAuth          bool
	ClientUser          string
	ClientPass          string
	Clienmaxidle        int
	Clienmaxperhost     int
	Clienmaxidleperhost int
	Clienidletimeout    time.Duration
	Clientimeout        time.Duration
	Monenabled          bool
	Monurl              string
	UpstreamsFile       string
	Accesslog           bool
	//Monuser       string
	//Monpass       string

}

var To = &confVars{
	HttpAddress:         "127.0.0.1:8080",
	Healtchecks:         20,
	ServerAuth:          false,
	BasicAuth:           false,
	JWTAuth:             false,
	ApiKeyAuth:          false,
	BasicCreds:          "",
	ClientAuth:          false,
	ClientUser:          "",
	ClientPass:          "",
	Monenabled:          false,
	Clienmaxidle:        100,
	Clienmaxperhost:     10,
	Clienmaxidleperhost: 10,
	Clienidletimeout:    90 * time.Second,
	Clientimeout:        time.Second * 10,
	Monurl:              "127.0.0.1:9191",
	UpstreamsFile:       "",
	Accesslog:           false,
	ApiKey:              os.Getenv("GAZANKEY"),
	JWTSecret:           []byte(os.Getenv("JWTSECRET")),
	//Monuser:       "",
	//Monpass:       "",

}

//var Authorized = make(map[string]string, 10)

func SetVarsik() {
	up := flag.String("up", "", "up")
	cfgFile := flag.String("config", "config.ini", "a string")
	flag.Parse()
	fmt.Println("Using :", *cfgFile, "as config file")

	cfg, err := ini.Load(*cfgFile)
	if err != nil {
		fmt.Printf("Fail To.read config file: %v", err)
		os.Exit(1)
	}
	To.UpstreamsFile = *up
	To.HttpAddress = cfg.Section("main").Key("listen").String()

	To.Healtchecks, _ = cfg.Section("main").Key("dispatchers").Int()

	To.ServerAuth, _ = cfg.Section("server").Key("serverauth").Bool()

	To.ClientAuth, _ = cfg.Section("client").Key("clientauth").Bool()
	To.ClientUser = cfg.Section("client").Key("clientuser").String()
	To.ClientPass = cfg.Section("client").Key("clientpass").String()

	To.Clienmaxidle, _ = cfg.Section("client").Key("maxidle").Int()
	To.Clienmaxperhost, _ = cfg.Section("client").Key("maxperhost").Int()
	To.Clienmaxidleperhost, _ = cfg.Section("client").Key("maxidleperhost").Int()
	tid, _ := cfg.Section("client").Key("idletimeout").Int()
	tdu, _ := cfg.Section("client").Key("timeout").Int()
	To.Clienidletimeout = time.Duration(tid)
	To.Clientimeout = time.Duration(tdu)

	To.Monenabled, _ = cfg.Section("monitoring").Key("enabled").Bool()
	To.Monurl = cfg.Section("monitoring").Key("url").String()
	//To.Monuser = cfg.Section("monitoring").Key("user").String()
	//To.Monpass = cfg.Section("monitoring").Key("pass").String()
	To.Accesslog = strtoBool(cfg.Section("server").Key("Accesslog").String())
	//Authorized["mon"] = To.Monuser + ":" + To.Monpass

	authtype := cfg.Section("server").Key("authtype").String()
	switch authtype {
	case "none":
		To.ServerAuth = false
		To.ApiKeyAuth = false
		To.BasicAuth = false
		To.JWTAuth = false
	case "apikey":
		if os.Getenv("GAZANKEY") == "" {
			log.Println("\n\n Api-Key authentication is enable but Key is not set \n Please set OS enviroment variable GAZANKEY to your api key\n")
			os.Exit(2)
		}
		To.ServerAuth = true
		To.ApiKeyAuth = true
		To.BasicAuth = false
		To.JWTAuth = false
	case "basic":
		if os.Getenv("BASICUSER") == "" && os.Getenv("BASICPASS") == "" {
			log.Println("\n\n Basic authentication is enable but user:password are not set \n Please set OS enviroment variable BASICUSER and  BASICPASS\n")
			os.Exit(2)
		}
		To.ServerAuth = true
		To.ApiKeyAuth = false
		To.BasicAuth = true
		To.JWTAuth = false
		To.BasicCreds = os.Getenv("BASICUSER") + ":" + os.Getenv("BASICPASS")
		//Authorized["server"] = os.Getenv("BASICUSER") + ":" + os.Getenv("BASICPASS")
	case "jwt":
		if os.Getenv("JWTSECRET") == "" {
			log.Println("\n\n JWT authentication is enable but SECRET is not set \n Please set OS enviroment variable JWTSECRET to your JWT secret\n")
			os.Exit(2)
		}
		To.ServerAuth = true
		To.ApiKeyAuth = false
		To.BasicAuth = false
		To.JWTAuth = true
	default:
		log.Println("Unknown authentication parameter")
		os.Exit(2)
	}
	fmt.Println("Authentication enabled:", authtype)
}
