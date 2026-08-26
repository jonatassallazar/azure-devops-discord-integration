package azuredevops

import (
	"fmt"
	"net/http"

	"azuredevops-notify/internal/notify"

	"github.com/gin-gonic/gin"
)

// PullRequestHandler handles Azure DevOps pull request webhooks and
// dispatches a notify.Message to every configured sink. It holds no
// per-request state, so one instance safely serves concurrent requests.
type PullRequestHandler struct {
	Dispatcher *notify.Dispatcher
	Avatars    *AvatarProxy
}

func (h *PullRequestHandler) CreatedPR(c *gin.Context) {
	var req AzureRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, err)
		return
	}

	msg := req.toPRMessage("Pull Request Criado", notify.LevelPending, h.Avatars)

	if err := h.Dispatcher.Send(c.Request.Context(), msg); err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, req)
}

func (h *PullRequestHandler) ReviewedPR(c *gin.Context) {
	var req AzureRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, err)
		return
	}

	if len(req.Resource.Reviewers) == 0 {
		c.JSON(http.StatusNoContent, gin.H{})
		return
	}

	level, title := h.processVotes(&req)
	if title == "" {
		c.JSON(http.StatusNoContent, gin.H{})
		return
	}

	msg := req.toPRMessage(fmt.Sprintf("Pull Request | %s", title), level, h.Avatars)

	if err := h.Dispatcher.Send(c.Request.Context(), msg); err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, req)
}

func (h *PullRequestHandler) StatusUpdatedPR(c *gin.Context) {
	var req AzureRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, err)
		return
	}

	var level notify.Level
	var title string

	switch req.Resource.Status {
	case "completed":
		level, title = notify.LevelCompleted, "Concluído"
	case "conflicts":
		level, title = notify.LevelFailure, "com Conflito"
	default:
		c.JSON(http.StatusNoContent, gin.H{})
		return
	}

	msg := req.toPRMessage(fmt.Sprintf("Pull Request %s", title), level, h.Avatars)

	if err := h.Dispatcher.Send(c.Request.Context(), msg); err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, req)
}

// processVotes maps reviewer votes to a notification level, giving any
// rejection precedence over waiting, and waiting precedence over approval.
func (h *PullRequestHandler) processVotes(req *AzureRequest) (notify.Level, string) {
	var approved, reproved, waiting bool

	for _, r := range req.Resource.Reviewers {
		switch r.Vote {
		case 10:
			approved = true
		case -10:
			reproved = true
		case -5:
			waiting = true
		}
	}

	switch {
	case reproved:
		return notify.LevelFailure, "Reprovado"
	case waiting:
		return notify.LevelWarning, "Aguardando Autor"
	case approved:
		return notify.LevelSuccess, "Aprovado"
	default:
		return 0, ""
	}
}
