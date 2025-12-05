package user

import (
	"encoding/json"
	"errors"

	"cozeva.com/vault/interfaces"
	"github.com/gin-gonic/gin"
)

func GetCurrentUser(c *gin.Context) (interfaces.UserPayload, error) {
	var currentUser interfaces.UserPayload

	currentUserAny, exists := c.Get("currentUser")
	if !exists {
		return currentUser, errors.New("current user not found in context")
	}

	currentUserAnyJson, err := json.Marshal(currentUserAny)
	if err != nil {
		return currentUser, err
	}

	err = json.Unmarshal(currentUserAnyJson, &currentUser)
	return currentUser, err
}
