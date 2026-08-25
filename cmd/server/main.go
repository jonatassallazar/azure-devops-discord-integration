package main

import (
	"log"

	"azuredevops-notify/internal/config"
	"azuredevops-notify/internal/httpserver"
)

func main() {
	var cfg config.Config

	if err := cfg.LoadEnvironment(); err != nil {
		log.Fatalln(err.Error())
	}

	s := httpserver.Server{Config: &cfg}

	if _, err := s.Init(); err != nil {
		log.Fatalln(err.Error())
	}
}
