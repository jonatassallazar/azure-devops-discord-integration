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

// ReleaseController manages the processing of Azure DevOps release events
// and sends notifications to Discord when release status changes occur.
//
// The controller receives webhooks from Azure DevOps, processes release information,
// and sends formatted messages to Discord based on the deployment status.
type ReleaseController struct {
	// ConfigServer contains server configuration, including Azure DevOps and Discord URLs,
	// as well as credentials needed for authentication.
	ConfigServer *config.ConfigServer
	
	// Response stores the Azure DevOps request containing information about the release
	// event (deployment status, release details, etc.).
	Response *models.AzureRequest
}

// ReleaseStatusReport is an HTTP handler that processes release status webhooks
// from Azure DevOps and sends notifications to Discord.
//
// The method performs the following operations:
//  1. Binds the JSON received in the request to the AzureRequest structure
//  2. Processes the release deployment status to determine the color and title of the notification
//  3. Converts the release data to Discord payload format
//  4. Sends the notification to the configured Discord webhook
//
// Returns:
//   - Status 200 (OK) with release data on success
//   - Status 204 (No Content) if the status is not mapped or does not require notification
//   - Status 400 (Bad Request) on processing error
func (p *ReleaseController) ReleaseStatusReport(c *gin.Context) {
	err := c.ShouldBindJSON(&p.Response)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"err": err,
		})
		return
	}

	color, title := p.processReleaseStatus()

	if title == "" {
		c.JSON(http.StatusNoContent, gin.H{})
		return
	}

	body := p.Response.ConvertToDiscordPayloadRelease(fmt.Sprintf("Release %s", title), color)

	json_data, err := json.Marshal(body)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"err": err,
		})
		return
	}

	_, err = http.Post(p.ConfigServer.DiscordEnvReleaseUrl, HEADER_APP_JSON, bytes.NewBuffer(json_data))
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"err": err,
		})
		return
	}

	c.JSON(http.StatusOK, p.Response)
}

// processReleaseStatus maps the Azure DevOps release deployment status to a color and title
// that will be used in the Discord notification.
//
// The method processes the deployment status from the release resource and maps it as follows:
//   - "succeeded" -> Green (GREEN) and title "Concluída"
//   - "failed" -> Red (RED) and title "Falhada"
//   - "stopped" -> Orange (ORANGE) and title "Interrompida"
//   - Any other status -> White (WHITE) and title with the unmapped status
//
// Returns:
//   - color: hexadecimal color code as int32 for the Discord embed
//   - title: title describing the release deployment status (empty string if status cannot be determined)
func (p *ReleaseController) processReleaseStatus() (int32, string) {
	switch p.Response.Resource.Deployment.DeploymentStatus {
	case "succeeded":
		return models.GREEN, "Concluída"
	case "failed":
		return models.RED, "Falhada"
	case "stopped":
		return models.ORANGE, "Interrompida"
	default:
		return models.WHITE, fmt.Sprintf("[Status não mapeado: %s]", p.Response.Resource.Deployment.DeploymentStatus)
	}
}
