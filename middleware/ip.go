package middleware

import (
	"net/http"

	"cozeva.com/vault/config"
	"cozeva.com/vault/pkg/validation"
	"github.com/gin-gonic/gin"
)

func CheckIP(c *gin.Context) {
	validIPRange := config.GetConfigValue("api.config.allowediprange")

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
