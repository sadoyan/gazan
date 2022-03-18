package mainfiles

import (
	"flag"
	"fmt"
	"gopkg.in/ini.v1"
	"os"
)

type confVars struct {
	httpAddress string
	//destinationURL   string
	dispatchersCount int
	serverAuth       bool
	serverUser       string
	serverPass       string
	clientAuth       bool
	clientUser       string
	clientPass       string
	internalQueue    bool
	rQueueName       string
	monenabled       bool
	monurl           string
	monuser          string
	monpass          string
	upstreamsFile    string
	//queue            chan map[string][]byte

}

var to = &confVars{
	httpAddress: "127.0.0.1:8080",
	//destinationURL:   "http://127.0.0.1:8000",
	dispatchersCount: 20,
	serverAuth:       false,
	serverUser:       "",
	serverPass:       "",
	clientAuth:       false,
	clientUser:       "",
	clientPass:       "",
	internalQueue:    false,
	rQueueName:       "oemetrics",
	monenabled:       false,
	monurl:           "127.0.0.1:9191",
	monuser:          "",
	monpass:          "",
	upstreamsFile:    "",
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
		fmt.Printf("Fail to read config file: %v", err)
		os.Exit(1)
	}
	to.upstreamsFile = *up
	to.httpAddress = cfg.Section("main").Key("listen").String()
	//to.destinationURL = cfg.Section("main").Key("remote").String()
	to.dispatchersCount, _ = cfg.Section("main").Key("dispatchers").Int()
	to.internalQueue, _ = cfg.Section("main").Key("internalqueue").Bool()
	//qs, _ := cfg.Section("main").Key("queuesize").Int()
	//to.queue = make(chan map[string][]byte, qs)

	to.serverAuth, _ = cfg.Section("server").Key("serverauth").Bool()
	to.serverUser = cfg.Section("server").Key("serveruser").String()
	to.serverPass = cfg.Section("server").Key("serverpass").String()

	to.clientAuth, _ = cfg.Section("client").Key("clientauth").Bool()
	to.clientUser = cfg.Section("client").Key("clientuser").String()
	to.clientPass = cfg.Section("client").Key("clientpass").String()

	to.monenabled, _ = cfg.Section("monitoring").Key("enabled").Bool()
	to.monurl = cfg.Section("monitoring").Key("url").String()
	to.monuser = cfg.Section("monitoring").Key("user").String()
	to.monpass = cfg.Section("monitoring").Key("pass").String()

	authorized["server"] = to.serverUser + ":" + to.serverPass
	authorized["mon"] = to.monuser + ":" + to.monpass
}
