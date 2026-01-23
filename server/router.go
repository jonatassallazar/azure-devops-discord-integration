package server

import (
	"discord-azure-integration/controllers"

	"github.com/gin-gonic/gin"
)

// SetupRouter configures and returns a Gin router with all webhook endpoints
// for Azure DevOps events.
//
// The method performs the following operations:
//  1. Sets the Gin framework mode based on configuration (if provided)
//  2. Creates a new Gin router with default middleware
//  3. Initializes controllers for pull requests, pipelines, and releases
//  4. Registers all POST endpoints for webhook events
//
// Registered routes:
//   - POST /pull-request/created - Handles pull request creation events
//   - POST /pull-request/review - Handles pull request review events
//   - POST /pull-request/status - Handles pull request status change events
//   - POST /pipeline/ - Handles pipeline status change events
//   - POST /release/ - Handles release status change events
//
// Returns:
//   - *gin.Engine: A configured Gin router ready to handle HTTP requests
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
