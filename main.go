package main

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.POST("/pull-request", func(c *gin.Context) {
		response, err := json.Marshal(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
		}

		c.JSON(http.StatusOK, response)
	})
	r.Run()
}
