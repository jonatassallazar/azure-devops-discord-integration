package httpserver

import (
	"azuredevops-notify/internal/config"

	"github.com/gin-gonic/gin"
)

type Server struct {
	Config *config.Config
}

func (s *Server) Init() (*gin.Engine, error) {
	r := s.SetupRouter()

	if err := r.Run(); err != nil {
		return nil, err
	}

	return r, nil
}
