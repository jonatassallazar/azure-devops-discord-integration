package server

import (
	"discord-azure-integration/controllers"

	"github.com/gin-gonic/gin"
)

func (s *Server) SetupRouter() *gin.Engine {
	if s.ConfigServer.GinMode != "" {
		gin.SetMode(s.ConfigServer.GinMode)
	}

	r := gin.Default()

	var p = controllers.PullRequestController{
		ConfigServer: s.ConfigServer,
	}

	r.POST(controllers.CREATED_ROUTE, p.CreatedPR)

	r.POST(controllers.REVIEW_ROUTE, p.ReviewedPR)

	r.POST(controllers.STATUS_ROUTE, p.StatusUpdatedPR)

	r.POST(controllers.PIPELINE_ROUTE, p.PipelineStatusReport)

	r.POST(controllers.RELEASE_ROUTE, p.ReleaseStatusReport)

	return r
}
