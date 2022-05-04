package main

import (
	"bytes"
	models "discord-azure-integration/Models"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.POST("/pull-request/created", func(c *gin.Context) {
		var res models.AzureRequest
		err := c.ShouldBindJSON(&res)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"err": err,
			})
		}

		body := res.ConvertToDiscordPayload("Pull Request Criado", 0, false)

		json_data, err := json.Marshal(body)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{})
		}

		_, err = http.Post("https://discord.com/api/webhooks/970868298713559101/DwjQ6AHT65e2Xrf0QwEbKwMpCXGRJu8mUyPegCCnyRShR7NYMP2Mi1i3rndrdAMQIKPy", "application/json", bytes.NewBuffer(json_data))

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{})
		}

		c.JSON(http.StatusOK, res)
	})

	r.POST("/pull-request/updated", func(c *gin.Context) {
		var res models.AzureRequest
		err := c.ShouldBindJSON(&res)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"err": err,
			})
		}

		var approved int8
		var reproved bool

		if len(res.Resource.Reviewers) < 1 {
			c.JSON(http.StatusBadRequest, gin.H{})
		} else {
			for _, i := range res.Resource.Reviewers {
				if i.Vote == 10 {
					approved += i.Vote
				} else if i.Vote == -10 {
					reproved = true
				}
			}
		}

		body := res.ConvertToDiscordPayload("Pull Request Atualizado", approved, reproved)

		json_data, err := json.Marshal(body)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{})
		}

		_, err = http.Post("https://discord.com/api/webhooks/970868298713559101/DwjQ6AHT65e2Xrf0QwEbKwMpCXGRJu8mUyPegCCnyRShR7NYMP2Mi1i3rndrdAMQIKPy", "application/json", bytes.NewBuffer(json_data))

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{})
		}

		c.JSON(http.StatusOK, res)
	})

	r.Run()
}
