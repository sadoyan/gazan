package main

import (
	"mainfiles"
	"utils"
)

func main() {
	go utils.GetSeed()
	//utils.GetHostsByDNS()
	mainfiles.RunServer()
}
