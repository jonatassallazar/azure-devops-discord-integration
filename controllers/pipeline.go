package controllers

import (
	"bytes"
	config "discord-azure-integration/config"
	models "discord-azure-integration/models"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PipelineController struct {
	ConfigServer *config.ConfigServer
	Response     *models.AzureRequest
}

func (p *PipelineController) PipelineStatusReport(c *gin.Context) {
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

	repository, err := p.getRepository()
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"err": err,
		})
		return
	}

	body := p.Response.ConvertToDiscordPayloadPipeline(fmt.Sprintf("Pipeline %s", title), color, repository)

	json_data, err := json.Marshal(body)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"err": err,
		})
		return
	}

	_, err = http.Post(p.ConfigServer.DiscordEnvPipelineUrl, HEADER_APP_JSON, bytes.NewBuffer(json_data))
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"err": err,
		})
		return
	}

	c.JSON(http.StatusOK, p.Response)
}

func (p *PipelineController) processStatus() (int32, string) {
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

func (p *PipelineController) getRepository() (models.AzureRepository, error) {
	repository_id := p.Response.Resource.TriggerInfo.TriggerRepository
	url := fmt.Sprintf("%s/%s/_apis/git/repositories/%s?api-version=5.1", p.ConfigServer.AzureOrganization, p.ConfigServer.AzureProject, repository_id)

	req, err := http.NewRequest(http.MethodGet, url, bytes.NewBuffer([]byte{}))
	if err != nil {
		return models.AzureRepository{}, err
	}

	authToken := base64.StdEncoding.EncodeToString([]byte(":" + p.ConfigServer.AzurePAT))

	req.Header.Add("Authorization", fmt.Sprintf("Basic %s", authToken))

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return models.AzureRepository{}, err
	}
	if response.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(response.Body)
		return models.AzureRepository{}, fmt.Errorf("azure returned error, status code: %d, body: %s", response.StatusCode, string(bodyBytes))

	}
	defer response.Body.Close()

	repository := models.AzureRepository{}
	err = json.NewDecoder(response.Body).Decode(&repository)
	if err != nil {
		return models.AzureRepository{}, err
	}

	return repository, nil
}
