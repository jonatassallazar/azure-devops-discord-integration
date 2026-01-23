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

// PipelineController manages the processing of Azure DevOps pipeline events
// and sends notifications to Discord when pipeline status changes occur.
//
// The controller receives webhooks from Azure DevOps, processes pipeline information,
// fetches repository details, and sends a formatted message to Discord.
type PipelineController struct {
	// ConfigServer contains server configuration, including Azure DevOps and Discord URLs,
	// as well as credentials needed for authentication.
	ConfigServer *config.ConfigServer
	
	// Response stores the Azure DevOps request containing information about the pipeline
	// event (status, repository, definition, etc.).
	Response *models.AzureRequest
}

// PipelineStatusReport is an HTTP handler that processes pipeline status webhooks
// from Azure DevOps and sends notifications to Discord.
//
// The method performs the following operations:
//  1. Binds the JSON received in the request to the AzureRequest structure
//  2. Processes the pipeline status to determine the color and title of the notification
//  3. Fetches detailed repository information from Azure DevOps
//  4. Converts the data to Discord payload format
//  5. Sends the notification to the configured Discord webhook
//
// Returns:
//   - Status 200 (OK) with pipeline data on success
//   - Status 204 (No Content) if the status is not mapped or does not require notification
//   - Status 400 (Bad Request) on processing error
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

// processStatus maps the Azure DevOps pipeline status to a color and title
// that will be used in the Discord notification.
//
// Supported statuses are:
//   - "succeeded" -> Green (GREEN) and title "Concluída"
//   - "failed" -> Red (RED) and title "Falhada"
//   - "stopped" -> Orange (ORANGE) and title "Interrompida"
//   - Any other status -> White (WHITE) and title with the unmapped status
//
// Returns:
//   - color: hexadecimal color code as int32 for the Discord embed
//   - title: title describing the pipeline status
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

// getRepository fetches detailed Git repository information from Azure DevOps
// using the Azure DevOps REST API.
//
// The method performs an authenticated GET request to the Azure DevOps API
// using the configured Personal Access Token (PAT). Authentication is done
// using Basic Auth with the token encoded in base64.
//
// Returns:
//   - AzureRepository: structure containing repository information (name, URL, etc.)
//   - error: error if the request fails, HTTP status is not 200, or JSON
//     cannot be decoded
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
