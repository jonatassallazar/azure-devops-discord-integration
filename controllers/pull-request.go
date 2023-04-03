package controllers

import (
	"bytes"
	config "discord-azure-integration/Config"
	models "discord-azure-integration/Models"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PullRequestController struct {
	ConfigServer *config.ConfigServer
	Response     *models.AzureRequest
}

func (p *PullRequestController) CreatedPR(c *gin.Context) {
	err := c.ShouldBindJSON(&p.Response)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"err": err,
		})
		return
	}

	body := p.Response.ConvertToDiscordPayload("Pull Request Criado", models.YELLOW)

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

	body := p.Response.ConvertToDiscordPayload(fmt.Sprintf("Pull Request | %s", title), color)

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
		c.JSON(http.StatusOK, gin.H{})
		return
	}

	body := p.Response.ConvertToDiscordPayload(fmt.Sprintf("Pull Request %s", title), color)

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
