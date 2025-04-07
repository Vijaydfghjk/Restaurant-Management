package middleware

import (
	"fmt"
	"net/http"
	tokenservice "restaurant-management/Token_service"

	"github.com/gin-gonic/gin"
)

func Authentication(c *gin.Context) {

	clientToken := c.Request.Header.Get("token")

	if clientToken == "" {

		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("No Authorization header provided")})
		c.Abort()
		return
	}

	cliams, err := tokenservice.Validatetoken(clientToken)

	if err != "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		c.Abort()
		return
	}

	c.Set("email", cliams.Email)
	c.Set("first_name", cliams.First_name)
	c.Set("last_name", cliams.Last_name)
	c.Set("uid", cliams.Uid)

	c.Next()

}
