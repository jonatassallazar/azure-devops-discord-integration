package controllers

import (
	"bytes"
	config "discord-azure-integration/config"
	models "discord-azure-integration/models"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PipelineController struct {
	ConfigServer *config.ConfigServer
	Response     *models.AzureRequest
}

func (p *PullRequestController) PipelineStatusReport(c *gin.Context) {
	err := c.ShouldBindJSON(&p.Response)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"err": err,
		})
		return
	}

	color, title := p.processStatus()

	if title == "" {
		c.JSON(http.StatusNoContent, gin.H{})
		return
	}

	body := p.Response.ConvertToDiscordPayloadPipeline(fmt.Sprintf("Pipeline %s", title), color)

	json_data, err := json.Marshal(body)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"err": err,
		})
		return
	}

	_, err = http.Post(p.ConfigServer.DiscordEnvBuildUrl, HEADER_APP_JSON, bytes.NewBuffer(json_data))
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"err": err,
		})
		return
	}

	c.JSON(http.StatusOK, p.Response)
}

func (p *PullRequestController) processStatus() (int32, string) {
	switch p.Response.Resource.Result {
	case "succeeded":
		return models.GREEN, "Concluída"
	case "failed":
		return models.RED, "Falhada"
	case "stopped":
		return models.ORANGE, "Interrompida"
	default:
		return models.WHITE, fmt.Sprintf("[Status não mapeado: %s]", p.Response.Resource.Result)
	}
}
