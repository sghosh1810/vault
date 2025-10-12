package middleware

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"shounak.me/configmanager/pkg/validation"
)

func CheckIP(c *gin.Context) {
	validIPRange := os.Getenv("api.allowediprange")

	if validIPRange == "" {
		c.Next()
		return
	}

	errMsg := validation.IsIPInRange(c.ClientIP(), validIPRange)

	if errMsg != nil {
		c.AbortWithStatusJSON(http.StatusForbidden, errMsg)
		return
	}

	c.Next()
}
