package main

import (
	"configs"
	"mainfiles"
	"utils"
)

func main() {
	//go utils.GetSeed()
	//utils.GetHostsByDNS()
	//go utils.GetHostsByHTTP()
	//utils.CheckJWTtoken()
	go utils.Valod(configs.To.Healtchecks)

	mainfiles.RunServer()
}
