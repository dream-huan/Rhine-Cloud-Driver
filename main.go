package main

import (
	"golandproject/Router"
	"golandproject/config"
)

func main() {
	config.ReadIni()
	Router.InitRouter()
}
