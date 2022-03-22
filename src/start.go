package main

import (
	"mainfiles"
	"utils"
)

func main() {
	//go utils.GetSeed()
	//utils.GetHostsByDNS()
	//go utils.GetHostsByHTTP()

	go utils.Valod(mainfiles.To.Healtchecks)
	mainfiles.RunServer()
}
