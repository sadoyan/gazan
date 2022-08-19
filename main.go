package main

import (
	"configs"
	"mainfiles"
	"utils"
)

func main() {
	//ankap.Play()
	//go utils.GetSeed()
	//utils.GetHostsByDNS()
	//go utils.GetHostsByHTTP()
	//utils.CheckJWTtoken()

	go utils.Valod(configs.To.Healtchecks)
	go utils.GetHostsbyDNS()
	mainfiles.RunServer()
}
