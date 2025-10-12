package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"shounak.me/configmanager/interfaces"
	sqlquery "shounak.me/configmanager/sql"
)

func MapSecretToProject(c *gin.Context) {
	var payload interfaces.MapOperationPayload

	if err := c.BindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "Failed to decode JSON",
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
