package server

import (
	config "discord-azure-integration/config"

	"github.com/gin-gonic/gin"
)

// Server represents the HTTP server that handles Azure DevOps webhook events
// and forwards notifications to Discord.
//
// The server uses the Gin web framework and contains configuration settings
// needed for routing and handling webhook requests.
type Server struct {
	// ConfigServer contains all configuration values including Azure DevOps
	// and Discord webhook URLs, authentication tokens, and server settings.
	ConfigServer *config.ConfigServer
}

// Init initializes and starts the HTTP server.
//
// The method performs the following operations:
//  1. Sets up the router with all webhook endpoints via SetupRouter()
//  2. Starts the HTTP server using Gin's default settings (listens on :8080)
//
// The server will block and listen for incoming webhook requests from Azure DevOps.
// The server runs until it is stopped or encounters an error.
//
// Returns:
//   - *gin.Engine: The configured Gin router (returned after server starts)
//   - error: Error if the server fails to start (e.g., port already in use)
//
// Note: This method blocks execution. To run the server in a non-blocking manner,
// consider using goroutines or calling r.Run() separately after SetupRouter().
func (s *Server) Init() (*gin.Engine, error) {
	r := s.SetupRouter()

	err := r.Run()
	if err != nil {
		return nil, err
	}

	return r, nil
}
