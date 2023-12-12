package main

import (
	config "discord-azure-integration/config"
	"discord-azure-integration/server"
	"log"
)

func main() {
	var c config.ConfigServer

	err := c.LoadEnvironment()
	if err != nil {
		log.Fatalln(err.Error())
	}

	var s = server.Server{
		ConfigServer: &c,
	}

	_, err = s.Init()
	if err != nil {
		log.Fatalln(err.Error())
	}
}
