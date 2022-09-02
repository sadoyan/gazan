package configs

import (
	"flag"
	"fmt"
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

type confVars struct {
	HttpAddress         string
	TLSEnabled          bool
	TLSAddress          string
	TLSCertFIle         string
	TLSPrivKey          string
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
	Confurl             string
	Dns                 bool
	DnsServer           []string
	DnsRecords          map[string]string
	//Monuser       string
	//Monpass       string

}

var To = &confVars{
	HttpAddress:         "127.0.0.1:8080",
	TLSEnabled:          false,
	TLSAddress:          "127.0.0.1:8443",
	TLSCertFIle:         "",
	TLSPrivKey:          "",
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
	Confurl:             "127.0.0.1:4141",
	UpstreamsFile:       "",
	Accesslog:           false,
	ApiKey:              os.Getenv("GAZANKEY"),
	JWTSecret:           []byte(os.Getenv("JWTSECRET")),
	Dns:                 false,
	DnsServer:           make([]string, 0),
	DnsRecords:          make(map[string]string, 0),
	//Monuser:       "",
	//Monpass:       "",

}

//var Authorized = make(map[string]string, 10)

/*
func SetVarsikulik() {
	up := flag.String("up", "", "up")
	cfgFile := flag.String("config", "config.iwwwni", "a string")
	flag.Parse()
	log.Println("Using :", *cfgFile, "as config file")

	cfg, err := ini.Load(*cfgFile)
	if err != nil {
		fmt.Printf("Fail To.read config file: %v", err)
		os.Exit(1)
	}
	To.UpstreamsFile = *up
	To.HttpAddress = cfg.Section("main").Key("listen").String()

	To.TLSEnabled = strtoBool(cfg.Section("tls").Key("enbaled").String())
	if To.TLSEnabled {
		To.TLSCertFIle = cfg.Section("tls").Key("certificate").String()
		To.TLSPrivKey = cfg.Section("tls").Key("privatekey").String()
		To.TLSAddress = cfg.Section("tls").Key("listen").String()
	}

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
	To.Confurl = cfg.Section("main").Key("confurl").String()
	//To.Monuser = cfg.Section("monitoring").Key("user").String()
	//To.Monpass = cfg.Section("monitoring").Key("pass").String()
	To.Accesslog = strtoBool(cfg.Section("server").Key("accesslog").String())
	//Authorized["mon"] = To.Monuser + ":" + To.Monpass
	To.Dns = strtoBool(cfg.Section("main").Key("dnsapi").String())

	xs := cfg.Section("main").Key("dnsserver").String()
	dd := strings.Split(strings.Replace(xs, " ", "", -1), ",")
	To.DnsServer = dd

	authtype := cfg.Section("server").Key("authtype").String()
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
	log.Println("Authentication enabled:", authtype)

}
*/

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
	//To.HttpAddress = fmt.Sprintf("%v", data["main"]["listen"])
	To.HttpAddress = data["main"]["listen"].(string)
	To.TLSEnabled = strtoBool(data["tls"]["enbaled"].(string))
	if To.TLSEnabled {
		To.TLSCertFIle = data["tls"]["certificate"].(string)
		To.TLSPrivKey = data["tls"]["privatekey"].(string)
		To.TLSAddress = data["tls"]["listen"].(string)
	}
	To.Healtchecks, _ = data["main"]["dispatchers"].(int)
	//To.ServerAuth = strtoBool(data["server"]["serverauth"].(string))
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
	To.Accesslog = strtoBool(data["server"]["accesslog"].(string))
	To.Dns = strtoBool(data["main"]["dnsapi"].(string))
	To.DnsServer = strings.Split(strings.Replace(data["main"]["dnsserver"].(string), " ", "", -1), ",")

	authtype := data["server"]["authtype"].(string)
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

	fmt.Println("-----------------")
	//for ku := range data["main"]["dnsnames"].([]interface{}) {
	//	print(data["main"]["dnsnames"][ku])
	//}
	//fmt.Printf("%T", data["main"]["dnsnames"])

	for _, item := range data["main"]["dnsnames"].([]interface{}) {
		fmt.Println(item)
		fmt.Printf("%T", item)
		fmt.Println()

		//for k, v := range item.(map[string]string) {
		//	fmt.Println("k:", k, "v:", v)
		//}

	}

	//fmt.Println(data["main"]["dnsnames"].(map[string]interface{}))
	fmt.Println("-----------------")

	// - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - //
}
