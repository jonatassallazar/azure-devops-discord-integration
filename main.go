package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.POST("/pull-request", func(c *gin.Context) {
		c.JSON(http.StatusOK, c.Request.Body)
	})
	r.Run()
}
