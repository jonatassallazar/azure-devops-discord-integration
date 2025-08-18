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

	var pr = controllers.PullRequestController{
		ConfigServer: s.ConfigServer,
	}

	var pp = controllers.PipelineController{
		ConfigServer: s.ConfigServer,
	}

	var re = controllers.ReleaseController{
		ConfigServer: s.ConfigServer,
	}

	r.POST(controllers.CREATED_ROUTE, pr.CreatedPR)

	r.POST(controllers.REVIEW_ROUTE, pr.ReviewedPR)

	r.POST(controllers.STATUS_ROUTE, pr.StatusUpdatedPR)

	r.POST(controllers.PIPELINE_ROUTE, pp.PipelineStatusReport)

	r.POST(controllers.RELEASE_ROUTE, re.ReleaseStatusReport)

	return r
}
