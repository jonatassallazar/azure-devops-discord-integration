package main

import (
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.POST("/pull-request", func(c *gin.Context) {
		res, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"err": err,
			})
		}

		c.Data(http.StatusOK, "test", res)
	})
	r.Run()
}
