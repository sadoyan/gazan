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

	go utils.Healtcheck(configs.To.Healtchecks)
	go utils.GetHostsbyDNS()
	go utils.PopulateUSers()
	mainfiles.RunServer()
}
