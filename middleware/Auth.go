package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/ashish9868/meracloud/lib"
	"github.com/ashish9868/meracloud/models"
	"github.com/gin-gonic/gin"
)

type AuthorizationHeader struct {
	Authorization string `header:"Authorization"`
}

func BearerAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		var auth AuthorizationHeader
		c.BindHeader(&auth)
		user := models.User{}
		token := strings.Replace(auth.Authorization, "Bearer ", "", 1)
		result := lib.DbInstance().First(&user, "token = ? and token_expiry > ?", token, time.Now())

		if result.Error != nil {
			print(time.Now().Location().String())
			print("Token not found")
			c.JSON(http.StatusUnauthorized, "Unauthorized")
			c.Abort()
			return
		}

		c.Set("user", user)
		c.Next()
	}
}
