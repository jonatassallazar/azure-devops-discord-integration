package azuredevops

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func respondError(c *gin.Context, err error) {
	fmt.Println(err)
	c.JSON(http.StatusBadRequest, gin.H{"err": err})
}
