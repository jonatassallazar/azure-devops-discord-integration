package main

import (
	config "discord-azure-integration/config"
	"discord-azure-integration/server"
	"log"
)

// main is the entry point of the Azure DevOps to Discord integration service.
//
// The function performs the following operations:
//  1. Creates a ConfigServer instance to hold configuration values
//  2. Loads environment variables from .env files using LoadEnvironment()
//  3. Creates a Server instance with the loaded configuration
//  4. Initializes and starts the HTTP server to listen for Azure DevOps webhooks
//
// The application will:
//   - Load configuration from environment variables (see config.ConfigServer for required variables)
//   - Start an HTTP server listening on port 8080 (default Gin port)
//   - Handle POST requests to webhook endpoints for pull requests, pipelines, and releases
//   - Forward notifications to Discord webhooks based on Azure DevOps events
//
// The server runs until it is stopped or encounters a fatal error.
// If configuration loading fails or the server fails to start, the application
// will log the error and exit with a non-zero status code.
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
