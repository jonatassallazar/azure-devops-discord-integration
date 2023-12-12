package server

import (
	config "discord-azure-integration/config"

	"github.com/gin-gonic/gin"
)

type Server struct {
	ConfigServer *config.ConfigServer
}

func (s *Server) Init() (*gin.Engine, error) {
	r := s.SetupRouter()

	err := r.Run()
	if err != nil {
		return nil, err
	}

	return r, nil
}
