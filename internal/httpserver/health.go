package httpserver

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RouteHealth is the liveness/readiness probe path. It lives here rather
// than in internal/azuredevops because it is a property of the server, not
// of any inbound event source.
const RouteHealth = "/health"

// healthCheck reports that the process is up and serving. The service is
// stateless and has no downstream dependency it needs at rest — sinks are
// only contacted while handling a webhook — so being able to answer is the
// whole readiness condition.
func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
