package configs

import (
	"flag"
	"gopkg.in/yaml.v3"
	"io/ioutil"
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

type DnsRecords []struct {
	Name    string `yaml:"name"`
	Address string `yaml:"address"`
}

type confVars struct {
	ImRunning           bool
	HttpAddress         string
	TLSEnabled          bool
	TLSAddress          string
	TLSCertFIle         string
	TLSPrivKey          string
	Healtchecks         int
	ServerAuth          bool
	Configauth          bool
	Configkey           string
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
	Confurl             string
	Dns                 bool
	DnsServer           []string
	DDD                 map[string]string
	DnsRecords          map[string]string
	//Monuser       string
	//Monpass       string
}

var To = &confVars{
	ImRunning:           false,
	HttpAddress:         "127.0.0.1:8080",
	TLSEnabled:          false,
	TLSAddress:          "127.0.0.1:8443",
	TLSCertFIle:         "",
	TLSPrivKey:          "",
	Healtchecks:         20,
	ServerAuth:          false,
	Configauth:          false,
	Configkey:           "",
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
	Confurl:             "127.0.0.1:4141",
	UpstreamsFile:       "",
	Accesslog:           false,
	ApiKey:              os.Getenv("GAZANKEY"),
	JWTSecret:           []byte(os.Getenv("JWTSECRET")),
	Dns:                 false,
	DnsServer:           make([]string, 0),
	DnsRecords:          make(map[string]string, 0),
}

//var Authorized = make(map[string]string, 10)

func SetVarsik() {
	up := flag.String("up", "", "up")
	cfgFile := flag.String("config", "config.yml", "a string")
	flag.Parse()
	log.Println("Using :", *cfgFile, "as config file")

	To.UpstreamsFile = *up

	yfile, err := ioutil.ReadFile(*cfgFile)
	if err != nil {
		log.Fatal(err)
	}

	data := make(map[interface{}]map[interface{}]interface{})

	err2 := yaml.Unmarshal(yfile, &data)

	if err2 != nil {
		log.Fatal(err2)
	}
	// - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - //
	To.HttpAddress = data["main"]["listen"].(string)
	To.TLSEnabled = strtoBool(data["tls"]["enabled"].(string))
	if To.TLSEnabled {
		To.TLSCertFIle = data["tls"]["certificate"].(string)
		To.TLSPrivKey = data["tls"]["privatekey"].(string)
		To.TLSAddress = data["tls"]["listen"].(string)
	}
	To.Healtchecks, _ = data["main"]["dispatchers"].(int)
	To.ClientAuth = strtoBool(data["client"]["clientauth"].(string))
	To.ClientUser = data["client"]["clientuser"].(string)
	To.ClientPass = data["client"]["clientpass"].(string)
	To.Clienmaxidle, _ = data["client"]["dispatchers"].(int)
	To.Clienmaxperhost, _ = data["client"]["maxperhost"].(int)
	To.Clienmaxidleperhost, _ = data["client"]["maxidleperhost"].(int)
	To.Clienidletimeout, _ = data["client"]["idletimeout"].(time.Duration)
	To.Clientimeout, _ = data["client"]["timeout"].(time.Duration)

	To.Monenabled = strtoBool(data["monitoring"]["enabled"].(string))
	To.Monurl = data["monitoring"]["url"].(string)
	To.Confurl = data["main"]["confurl"].(string)
	To.Accesslog = strtoBool(data["main"]["accesslog"].(string))

	To.Configauth = strtoBool(data["main"]["configauth"].(string))
	switch To.Configauth {
	case true:
		if os.Getenv("CONFIGKEY") == "" {
			log.Print("\n\n\n Authentication for configuration is enable but Key is not set \n Please set OS enviroment variable CONFIGKEY to your api key\n\n")
			os.Exit(2)
		} else {
			To.Configkey = os.Getenv("CONFIGKEY")
		}
	}
	authtype := data["main"]["authtype"].(string)
	switch authtype {
	case "none":
		To.ServerAuth = false
		To.ApiKeyAuth = false
		To.BasicAuth = false
		To.JWTAuth = false
	case "apikey":
		if os.Getenv("GAZANKEY") == "" {
			log.Print("\n\n\n Api-Key authentication is enable but Key is not set \n Please set OS enviroment variable GAZANKEY to your api key\n\n")
			os.Exit(2)
		}
		To.ServerAuth = true
		To.ApiKeyAuth = true
		To.BasicAuth = false
		To.JWTAuth = false
	case "basic":
		if os.Getenv("BASICUSER") == "" && os.Getenv("BASICPASS") == "" {
			log.Print("\n\n\n Basic authentication is enable but user:password are not set \n Please set OS enviroment variable BASICUSER and  BASICPASS\n\n")
			os.Exit(2)
		}
		To.ServerAuth = true
		To.ApiKeyAuth = false
		To.BasicAuth = true
		To.JWTAuth = false
		To.BasicCreds = os.Getenv("BASICUSER") + ":" + os.Getenv("BASICPASS")
	case "jwt":
		if os.Getenv("JWTSECRET") == "" {
			log.Print("\n\n\n JWT authentication is enable but SECRET is not set \n Please set OS enviroment variable JWTSECRET to your JWT secret\n\n")
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

	//To.DnsServer = strings.Split(strings.Replace(data["main"]["dnsserver"].(string), " ", "", -1), ",")
	To.Dns = strtoBool(data["dnsconfig"]["enabled"].(string))
	for _, dnssrv := range data["dnsconfig"]["dnsservers"].([]interface{}) {
		To.DnsServer = append(To.DnsServer, dnssrv.(string))
	}
	for _, vvv := range data["dnsconfig"]["srvrecords"].([]interface{}) {
		To.DnsRecords[vvv.(map[string]interface{})["name"].(string)] = vvv.(map[string]interface{})["address"].(string)
	}
}
