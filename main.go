package main

import (
	"bytes"
	config "discord-azure-integration/Config"
	models "discord-azure-integration/Models"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	var configsUrls config.ConfigUrls

	config.LoadEnvironment(&configsUrls)

	r.POST("/pull-request/created", func(c *gin.Context) {
		var res models.AzureRequest
		err := c.ShouldBindJSON(&res)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"err": err,
			})
			return
		}

		fmt.Println(res)

		body := res.ConvertToDiscordPayload("Pull Request Criado", models.YELLOW)

		json_data, err := json.Marshal(body)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{})
			return
		}

		_, err = http.Post(configsUrls.DiscordEnvPRUrl, "application/json", bytes.NewBuffer(json_data))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{})
			return
		}

		c.JSON(http.StatusOK, res)
	})

	r.POST("/pull-request/review", func(c *gin.Context) {
		var res models.AzureRequest
		err := c.ShouldBindJSON(&res)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"err": err,
			})
			return
		}

		var approved int8
		var reproved bool

		if len(res.Resource.Reviewers) < 2 {
			c.JSON(http.StatusOK, gin.H{})
			return
		} else {
			for _, i := range res.Resource.Reviewers {
				if i.Vote == 10 {
					approved += i.Vote
				} else if i.Vote == -10 {
					reproved = true
				}
			}
		}

		var color int32
		var title string

		if approved >= 20 {
			color = models.GREEN
			title = "Aprovado"
		} else if reproved {
			color = models.RED
			title = "Reprovado"
		} else {
			c.JSON(http.StatusOK, gin.H{})
			return
		}

		body := res.ConvertToDiscordPayload(fmt.Sprintf("Pull Request %s", title), color)

		json_data, err := json.Marshal(body)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{})
			return
		}

		_, err = http.Post(configsUrls.DiscordEnvPRUrl, "application/json", bytes.NewBuffer(json_data))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{})
			return
		}

		c.JSON(http.StatusOK, res)
	})

	r.POST("/pull-request/status", func(c *gin.Context) {
		var res models.AzureRequest
		err := c.ShouldBindJSON(&res)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"err": err,
			})
			return
		}

		var color int32
		var title string

		if res.Resource.Status == "succeeded" {
			color = models.BLURPLE
			title = "Concluído"
		} else if res.Resource.Status == "conflicts" {
			color = models.RED
			title = "com Conflito"
		} else {
			c.JSON(http.StatusBadRequest, gin.H{})
			return
		}

		body := res.ConvertToDiscordPayload(fmt.Sprintf("Pull Request %s", title), color)

		json_data, err := json.Marshal(body)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{})
			return
		}

		_, err = http.Post(configsUrls.DiscordEnvPRUrl, "application/json", bytes.NewBuffer(json_data))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{})
			return
		}

		c.JSON(http.StatusOK, res)
	})

	r.POST("/build/completed", func(c *gin.Context) {
		var res models.AzureRequest
		err := c.ShouldBindJSON(&res)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"err": err,
			})
			return
		}

		body := res.ConvertToDiscordPayload("Deploy realizado com Sucesso", models.BLURPLE)

		json_data, err := json.Marshal(body)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{})
			return
		}

		_, err = http.Post(configsUrls.DiscordEnvBuildUrl, "application/json", bytes.NewBuffer(json_data))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{})
			return
		}

		c.JSON(http.StatusOK, res)
	})

	r.Run()
}
