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

// PullRequestController manages the processing of Azure DevOps pull request events
// and sends notifications to Discord when pull request events occur.
//
// The controller receives webhooks from Azure DevOps for pull request events
// (creation, reviews, status changes), processes the information, and sends
// formatted messages to Discord.
type PullRequestController struct {
	// ConfigServer contains server configuration, including Azure DevOps and Discord URLs,
	// as well as credentials needed for authentication.
	ConfigServer *config.ConfigServer
	
	// Response stores the Azure DevOps request containing information about the pull request
	// event (status, reviewers, votes, etc.).
	Response *models.AzureRequest
}

// CreatedPR is an HTTP handler that processes pull request created events
// from Azure DevOps and sends notifications to Discord.
//
// The method performs the following operations:
//  1. Binds the JSON received in the request to the AzureRequest structure
//  2. Converts the pull request data to Discord payload format with a yellow color
//  3. Sends the notification to the configured Discord webhook
//
// Returns:
//   - Status 200 (OK) with pull request data on success
//   - Status 400 (Bad Request) on processing error
func (p *PullRequestController) CreatedPR(c *gin.Context) {
	err := c.ShouldBindJSON(&p.Response)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"err": err,
		})
		return
	}

	body := p.Response.ConvertToDiscordPayloadPR("Pull Request Criado", models.YELLOW)

	json_data, err := json.Marshal(body)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"err": err,
		})
		return
	}

	_, err = http.Post(p.ConfigServer.DiscordEnvPRUrl, HEADER_APP_JSON, bytes.NewBuffer(json_data))
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"err": err,
		})
		return
	}

	c.JSON(http.StatusOK, p.Response)
}

// ReviewedPR is an HTTP handler that processes pull request review events
// from Azure DevOps and sends notifications to Discord.
//
// The method performs the following operations:
//  1. Binds the JSON received in the request to the AzureRequest structure
//  2. Validates that reviewers exist in the request
//  3. Processes reviewer votes to determine the review status (approved, rejected, waiting)
//  4. Converts the review data to Discord payload format with appropriate color
//  5. Sends the notification to the configured Discord webhook
//
// Returns:
//   - Status 200 (OK) with pull request data on success
//   - Status 204 (No Content) if no reviewers exist or status cannot be determined
//   - Status 400 (Bad Request) on processing error
func (p *PullRequestController) ReviewedPR(c *gin.Context) {
	err := c.ShouldBindJSON(&p.Response)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"err": err,
		})
		return
	}

	if len(p.Response.Resource.Reviewers) == 0 {
		c.JSON(http.StatusNoContent, gin.H{})
		return
	}

	color, title := p.processVotes()

	if title == "" {
		c.JSON(http.StatusNoContent, gin.H{})
		return
	}

	body := p.Response.ConvertToDiscordPayloadPR(fmt.Sprintf("Pull Request | %s", title), color)

	json_data, err := json.Marshal(body)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"err": err,
		})
		return
	}

	_, err = http.Post(p.ConfigServer.DiscordEnvPRUrl, HEADER_APP_JSON, bytes.NewBuffer(json_data))
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"err": err,
		})
		return
	}

	c.JSON(http.StatusOK, p.Response)
}

// StatusUpdatedPR is an HTTP handler that processes pull request status change events
// from Azure DevOps and sends notifications to Discord.
//
// The method handles the following status changes:
//   - "completed" -> Blurple color with "Concluído" title
//   - "conflicts" -> Red color with "com Conflito" title
//   - Any other status -> Returns 204 (No Content) without notification
//
// The method performs the following operations:
//  1. Binds the JSON received in the request to the AzureRequest structure
//  2. Determines the color and title based on the pull request status
//  3. Converts the status data to Discord payload format
//  4. Sends the notification to the configured Discord webhook
//
// Returns:
//   - Status 200 (OK) with pull request data on success
//   - Status 204 (No Content) if the status is not handled or does not require notification
//   - Status 400 (Bad Request) on processing error
func (p *PullRequestController) StatusUpdatedPR(c *gin.Context) {
	err := c.ShouldBindJSON(&p.Response)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"err": err,
		})
		return
	}

	var color int32
	var title string

	if p.Response.Resource.Status == "completed" {
		color = models.BLURPLE
		title = "Concluído"
	} else if p.Response.Resource.Status == "conflicts" {
		color = models.RED
		title = "com Conflito"
	} else {
		c.JSON(http.StatusNoContent, gin.H{})
		return
	}

	body := p.Response.ConvertToDiscordPayloadPR(fmt.Sprintf("Pull Request %s", title), color)

	json_data, err := json.Marshal(body)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"err": err,
		})
		return
	}

	_, err = http.Post(p.ConfigServer.DiscordEnvPRUrl, HEADER_APP_JSON, bytes.NewBuffer(json_data))
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"err": err,
		})
		return
	}

	c.JSON(http.StatusOK, p.Response)
}

// processVotes analyzes reviewer votes from the pull request and determines
// the overall review status and appropriate notification color.
//
// The method processes votes from all reviewers and determines the status based on:
//   - Vote 10: Approved (approves the pull request)
//   - Vote -10: Rejected (rejects the pull request)
//   - Vote -5: Waiting for author (requests changes)
//
// Priority order for status determination:
//  1. If any reviewer rejected (-10) -> Red color, "Reprovado" title
//  2. If any reviewer is waiting (-5) -> Orange color, "Aguardando Autor" title
//  3. If any reviewer approved (10) -> Green color, "Aprovado" title
//
// Returns:
//   - color: hexadecimal color code as int32 for the Discord embed
//   - title: title describing the review status (empty string if no status can be determined)
func (p *PullRequestController) processVotes() (int32, string) {
	var approved bool
	var reproved bool
	var waiting bool

	for _, i := range p.Response.Resource.Reviewers {
		switch i.Vote {
		case 10:
			approved = true
		case -10:
			reproved = true
		case -5:
			waiting = true
		}
	}

	var color int32
	var title string

	if reproved {
		color = models.RED
		title = "Reprovado"
	} else if waiting {
		color = models.ORANGE
		title = "Aguardando Autor"
	} else if approved {
		color = models.GREEN
		title = "Aprovado"
	}

	return color, title
}
