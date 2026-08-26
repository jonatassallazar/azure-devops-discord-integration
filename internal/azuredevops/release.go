package azuredevops

import (
	"fmt"
	"net/http"

	"azuredevops-notify/internal/notify"

	"github.com/gin-gonic/gin"
)

// ReleaseHandler handles Azure DevOps release deployment webhooks. It holds
// no per-request state, so one instance safely serves concurrent requests.
type ReleaseHandler struct {
	Dispatcher *notify.Dispatcher
}

func (h *ReleaseHandler) ReleaseStatusReport(c *gin.Context) {
	var req AzureRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, err)
		return
	}

	level, title := h.processReleaseStatus(&req)
	if title == "" {
		c.JSON(http.StatusNoContent, gin.H{})
		return
	}

	msg := req.toReleaseMessage(fmt.Sprintf("Release %s", title), level)

	if err := h.Dispatcher.Send(c.Request.Context(), msg); err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, req)
}

func (h *ReleaseHandler) processReleaseStatus(req *AzureRequest) (notify.Level, string) {
	switch req.Resource.Deployment.DeploymentStatus {
	case "succeeded":
		return notify.LevelSuccess, "Concluída"
	case "failed":
		return notify.LevelFailure, "Falhada"
	case "stopped":
		return notify.LevelWarning, "Interrompida"
	default:
		return notify.LevelUnmapped, fmt.Sprintf("[Status não mapeado: %s]", req.Resource.Deployment.DeploymentStatus)
	}
}
