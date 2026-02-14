package api

import (
	"fmt"
	"net/http"

	"cozeva.com/vault/interfaces"
	"cozeva.com/vault/pkg/access"
	"cozeva.com/vault/pkg/user"
	sqlquery "cozeva.com/vault/sql"
	"github.com/gin-gonic/gin"
)

func MapSecretToProject(c *gin.Context) {
	var payload interfaces.MapOperationPayload

	if err := c.BindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "Missing required parameter: project_id and secret_id",
		})
		return
	}

	currentUser, err := user.GetCurrentUser(c)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "Failed to parse current user info.",
		})
		return
	}

	projectAccess, err := access.GetAccessForProject(payload.ProjectID, payload.WorkspaceID, currentUser.Uid)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "Failed to get access for project.",
		})
		return
	}

	secretAccess, err := access.GetAccessForSecret(payload.SecretID, currentUser.Uid)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "Failed to get access for secret.",
		})
		return
	}

	if !projectAccess.HasWriteAccess || !secretAccess.HasWriteAccess {
		c.JSON(http.StatusForbidden, gin.H{
			"status":  "fail",
			"message": "Access denied to this resource",
		})
		return
	}

	db, err := sqlquery.GetSqlInstance()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "Failed to connect to db.",
		})
		return
	}
	defer db.Close()

	_, err = db.Exec(sqlquery.MapSecretToProject, payload.ProjectID, payload.SecretID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "Failed to map secret with project.",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": fmt.Sprintf("Successfully added secret with id  %d to project with id %d", payload.SecretID, payload.ProjectID),
	})
}
