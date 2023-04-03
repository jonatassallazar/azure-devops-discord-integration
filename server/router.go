package server

import (
	"discord-azure-integration/controllers"

	"github.com/gin-gonic/gin"
)

func (s *Server) SetupRouter() *gin.Engine {
	r := gin.Default()

	var p = controllers.PullRequestController{
		ConfigServer: s.ConfigServer,
	}

	r.POST("/pull-request/created", p.CreatedPR)

	r.POST("/pull-request/review", p.ReviewedPR)

	r.POST("/pull-request/status", p.StatusUpdatedPR)

	// r.POST("/build/completed")

	return r
}
