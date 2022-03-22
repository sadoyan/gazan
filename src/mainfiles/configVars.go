package mainfiles

import (
	"flag"
	"fmt"
	"gopkg.in/ini.v1"
	"os"
	"strings"
)

func strtoBool(in string) bool {
	out := false
	if strings.ToLower(in) == "on" || strings.ToLower(in) == "yes" {
		out = true
	}
	return out
}

type confVars struct {
	httpAddress string
	//destinationURL   string
	Healtchecks   int
	serverAuth    bool
	serverUser    string
	serverPass    string
	clientAuth    bool
	clientUser    string
	clientPass    string
	internalQueue bool
	rQueueName    string
	monenabled    bool
	monurl        string
	monuser       string
	monpass       string
	upstreamsFile string
	accesslog     bool
	//queue            chan map[string][]byte

}

var To = &confVars{
	httpAddress: "127.0.0.1:8080",
	//destinationURL:   "http://127.0.0.1:8000",
	Healtchecks:   20,
	serverAuth:    false,
	serverUser:    "",
	serverPass:    "",
	clientAuth:    false,
	clientUser:    "",
	clientPass:    "",
	internalQueue: false,
	rQueueName:    "oemetrics",
	monenabled:    false,
	monurl:        "127.0.0.1:9191",
	monuser:       "",
	monpass:       "",
	upstreamsFile: "",
	accesslog:     false,
	//queue:            make(chan map[string][]byte, 5000000),

}

var authorized = make(map[string]string, 10)

func setVarsik() {
	up := flag.String("up", "", "up")
	cfgFile := flag.String("config", "config.ini", "a string")
	flag.Parse()
	fmt.Println("Using :", *cfgFile, "as config file")

	cfg, err := ini.Load(*cfgFile)
	if err != nil {
		fmt.Printf("Fail To.read config file: %v", err)
		os.Exit(1)
	}
	To.upstreamsFile = *up
	To.httpAddress = cfg.Section("main").Key("listen").String()

	To.Healtchecks, _ = cfg.Section("main").Key("dispatchers").Int()
	To.internalQueue, _ = cfg.Section("main").Key("internalqueue").Bool()

	To.serverAuth, _ = cfg.Section("server").Key("serverauth").Bool()
	To.serverUser = cfg.Section("server").Key("serveruser").String()
	To.serverPass = cfg.Section("server").Key("serverpass").String()

	To.clientAuth, _ = cfg.Section("client").Key("clientauth").Bool()
	To.clientUser = cfg.Section("client").Key("clientuser").String()
	To.clientPass = cfg.Section("client").Key("clientpass").String()

	To.monenabled, _ = cfg.Section("monitoring").Key("enabled").Bool()
	To.monurl = cfg.Section("monitoring").Key("url").String()
	To.monuser = cfg.Section("monitoring").Key("user").String()
	To.monpass = cfg.Section("monitoring").Key("pass").String()
	To.accesslog = strtoBool(cfg.Section("server").Key("accesslog").String())
	authorized["server"] = To.serverUser + ":" + To.serverPass
	authorized["mon"] = To.monuser + ":" + To.monpass
}
