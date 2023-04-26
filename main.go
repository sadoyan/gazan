package main

import (
	"configs"
	"mainfiles"
	"utils"
)

func main() {
	go utils.Healtcheck(configs.To.Healtchecks)
	go utils.GetHostsbyDNS()
	utils.PopulateUSers()
	mainfiles.RunServer()
}
