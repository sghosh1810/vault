package middleware

import (
	"github.com/gin-gonic/gin"
)

func CheckAuth(c *gin.Context) {
	c.Next()

	// authHeader := c.GetHeader("Authorization")

	// claims, errMsg := security.VerifyJWT(authHeader)

	// if errMsg != nil {
	// 	c.AbortWithStatusJSON(http.StatusForbidden, errMsg)
	// 	return
	// }

	// c.Set("currentUser", claims)
	// c.Next()
}
