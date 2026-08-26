package azuredevops

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func respondError(c *gin.Context, err error) {
	respondErrorStatus(c, http.StatusBadRequest, err)
}

// respondErrorStatus is respondError for the cases that are not the
// caller's fault - an upstream Azure DevOps fetch that failed is a 502, not
// a bad request.
func respondErrorStatus(c *gin.Context, status int, err error) {
	fmt.Println(err)
	c.JSON(status, gin.H{"err": err})
}
