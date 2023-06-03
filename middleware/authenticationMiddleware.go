package middleware

import (
	"fmt"
	"net/http"

	"github.com/dev-hack95/jwt-golang/helpers"
	"github.com/gin-gonic/gin"
)

func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		clienToken := c.Request.Header.Get("token")
		if clienToken == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Header Not provided")})
			c.Abort()
			return
		}

		claims, err := helpers.ValidateToken(clienToken)

		if err != "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			c.Abort()
			return
		}

		c.Set("email", claims.Email)
		c.Set("firstname", claims.Firstname)
		c.Set("lastname", claims.Lastname)
		c.Set("uid", claims.Uid)
		c.Set("usertype", claims.UserType)
	}
}
