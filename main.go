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
	r.POST("/pull-request", func(c *gin.Context) {
		var res models.AzureRequest
		err := c.ShouldBindJSON(&res)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"err": err,
			})
		}

		Embeds := []models.Embeds{
			{
				Author:      models.Author{Name: res.Resource.LastChangedBy.DisplayName, Url: "", IconUrl: ""},
				Title:       "Pull Request Created",
				Url:         res.Resource.Url,
				Description: "",
				Color:       15258703,
				Fields:      []models.Field{{Name: "PR Aberto", Value: res.DetailedMessage.Text}},
				Thumbnail:   models.URL{Url: "https://upload.wikimedia.org/wikipedia/commons/3/38/4-Nature-Wallpapers-2014-1_ukaavUI.jpg"},
				Image:       models.URL{Url: ""},
				Footer:      models.Footer{Text: "", IconUrl: ""},
			},
		}

		body := models.DiscordPayload{
			Username:  "TestAzurePayload",
			AvatarUrl: "https://i.imgur.com/4M34hi2.png",
			Content:   res.Message.Text,
			Embeds:    Embeds,
		}

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
