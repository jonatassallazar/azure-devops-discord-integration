package main

import (
	config "discord-azure-integration/Config"
	"discord-azure-integration/server"
	"log"
)

func main() {
	var configsUrls config.ConfigUrls

	err := config.LoadEnvironment(&configsUrls)
	if err != nil {
		log.Fatalln(err.Error())
	}

	var s = server.Server{
		ConfigUrls: &configsUrls,
	}

	_, err = s.Init()
	if err != nil {
		log.Fatalln(err.Error())
	}
}
