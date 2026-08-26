package azuredevops

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"azuredevops-notify/internal/config"
	"azuredevops-notify/internal/notify"

	"github.com/gin-gonic/gin"
)

// PipelineHandler handles Azure DevOps pipeline (build) webhooks. It holds
// no per-request state, so one instance safely serves concurrent requests.
type PipelineHandler struct {
	Dispatcher *notify.Dispatcher
	Azure      config.AzureConfig
	Avatars    *AvatarProxy
}

func (h *PipelineHandler) PipelineStatusReport(c *gin.Context) {
	var req AzureRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, err)
		return
	}

	level, title := h.processStatus(&req)
	if title == "" {
		c.JSON(http.StatusNoContent, gin.H{})
		return
	}

	repository, err := h.getRepository(c.Request.Context(), &req)
	if err != nil {
		respondError(c, err)
		return
	}

	msg := req.toPipelineMessage(fmt.Sprintf("Pipeline %s", title), level, repository, h.Avatars)

	if err := h.Dispatcher.Send(c.Request.Context(), msg); err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, req)
}

func (h *PipelineHandler) processStatus(req *AzureRequest) (notify.Level, string) {
	switch req.Resource.Result {
	case "succeeded":
		return notify.LevelSuccess, "Concluída"
	case "failed":
		return notify.LevelFailure, "Falhada"
	case "stopped":
		return notify.LevelWarning, "Interrompida"
	default:
		return notify.LevelUnmapped, fmt.Sprintf("[Status não mapeado: %s]", req.Resource.Result)
	}
}

func (h *PipelineHandler) getRepository(ctx context.Context, req *AzureRequest) (AzureRepository, error) {
	repositoryID := req.Resource.TriggerInfo.TriggerRepository
	url := fmt.Sprintf("%s/%s/_apis/git/repositories/%s?api-version=5.1", h.Azure.Organization, h.Azure.Project, repositoryID)

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, url, bytes.NewBuffer([]byte{}))
	if err != nil {
		return AzureRepository{}, err
	}

	authToken := base64.StdEncoding.EncodeToString([]byte(":" + h.Azure.PATToken))
	httpReq.Header.Add("Authorization", fmt.Sprintf("Basic %s", authToken))

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return AzureRepository{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return AzureRepository{}, fmt.Errorf("azure returned error, status code: %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	var repository AzureRepository
	if err := json.NewDecoder(resp.Body).Decode(&repository); err != nil {
		return AzureRepository{}, err
	}

	return repository, nil
}
