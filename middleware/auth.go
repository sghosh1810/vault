package middleware

import (
	"net/http"

	"cozeva.com/vault/pkg/security"
	"github.com/gin-gonic/gin"
)

func CheckAuth(c *gin.Context) {

	authHeader := c.GetHeader("Authorization")

	claims, errMsg := security.VerifyJWT(authHeader)

	if errMsg != nil {
		c.AbortWithStatusJSON(http.StatusForbidden, errMsg)
		return
	}

	c.Set("currentUser", claims)
	c.Next()
}
