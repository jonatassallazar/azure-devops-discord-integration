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
}

func (p *PullRequestController) CreatedPR(c *gin.Context) {
	var res models.AzureRequest
	err := c.ShouldBindJSON(&res)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"err": err,
		})
		return
	}

	body := res.ConvertToDiscordPayload("Pull Request Criado", models.YELLOW)

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

	c.JSON(http.StatusOK, res)
}

func (p *PullRequestController) ReviewedPR(c *gin.Context) {
	var res models.AzureRequest
	err := c.ShouldBindJSON(&res)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"err": err,
		})
		return
	}

	var approved int8
	var reproved bool
	var waiting bool

	if len(res.Resource.Reviewers) == 0 {
		c.JSON(http.StatusOK, gin.H{})
		return
	} else {
		for _, i := range res.Resource.Reviewers {
			if i.Vote == 10 {
				approved += i.Vote
			} else if i.Vote == -10 {
				reproved = true
			} else if i.Vote == -5 {
				waiting = true
			}
		}
	}

	var color int32
	var title string

	if approved >= 10 {
		color = models.GREEN
		title = "Aprovado"
	} else if reproved {
		color = models.RED
		title = "Reprovado"
	} else if waiting {
		color = models.ORANGE
		title = "Aguardando Autor"
	} else {
		c.JSON(http.StatusNoContent, gin.H{})
		return
	}

	body := res.ConvertToDiscordPayload(fmt.Sprintf("Pull Request | %s", title), color)

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

	c.JSON(http.StatusOK, res)
}

func (p *PullRequestController) StatusUpdatedPR(c *gin.Context) {
	var res models.AzureRequest
	err := c.ShouldBindJSON(&res)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"err": err,
		})
		return
	}

	var color int32
	var title string

	if res.Resource.Status == "completed" {
		color = models.BLURPLE
		title = "Concluído"
	} else if res.Resource.Status == "conflicts" {
		color = models.RED
		title = "com Conflito"
	} else {
		c.JSON(http.StatusOK, gin.H{})
		return
	}

	body := res.ConvertToDiscordPayload(fmt.Sprintf("Pull Request %s", title), color)

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

	c.JSON(http.StatusOK, res)
}
