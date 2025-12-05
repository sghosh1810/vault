package api

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"cozeva.com/vault/interfaces"
	"cozeva.com/vault/pkg/access"
	"cozeva.com/vault/pkg/security"
	"cozeva.com/vault/pkg/user"
	sqlquery "cozeva.com/vault/sql"
	"github.com/gin-gonic/gin"
)

func SecretCreate(c *gin.Context) {
	var payload interfaces.SecretCreatePayload

	if err := c.BindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "Missing required parameter: name, value and environment",
		})
		return
	}

	if payload.SecretName == "" || payload.SecretValue == "" || payload.SecretEnvironment == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "Secret name, value and environment are required",
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

	db, err := sqlquery.GetSqlInstance()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "fail",
			"message": "Failed to connect to db.",
		})
		return
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "fail",
			"message": "Failed to start transaction",
		})
		return
	}

	res, err := tx.Exec(sqlquery.SecretInsertQuery, payload.SecretName)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "Failed to insert secret metadata",
		})
		return
	}
	secretId, _ := res.LastInsertId()

	encryptedData, err := security.EncryptSecret(payload.SecretValue)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "Failed to encrypt secret.",
		})
		return
	}

	_, err = tx.Exec(sqlquery.SecretVersionInsertQuery, secretId, 1, encryptedData, payload.SecretEnvironment)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "Failed to insert secret version",
		})
		return
	}

	_, err = tx.Exec(sqlquery.InsertUserSecretAccessQuery, currentUser.Uid, secretId, 1, 1, 1)

	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "Failed to create new secret with name " + payload.SecretName,
		})
		return
	}

	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "fail",
			"message": "Transaction commit failed",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": fmt.Sprintf("Created secret %s with version 1", payload.SecretName),
		"id":      secretId,
	})
}

func SecretUpdate(c *gin.Context) {
	var payload interfaces.SecretUpdatePayload

	currentUser, err := user.GetCurrentUser(c)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "Failed to parse current user info.",
		})
		return
	}

	if err := c.BindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "Missing required parameter: id and value",
		})
		return
	}

	if payload.SecretID == 0 || payload.SecretValue == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "Secret ID and value are required",
		})
		return
	}

	projectAccess, err := access.GetAccessForSecret(payload.SecretID, currentUser.Uid)

	if err != nil || !projectAccess.HasWriteAccess {
		c.JSON(http.StatusForbidden, gin.H{
			"status":  "fail",
			"message": "Access denied to this resource",
		})
		return
	}

	db, err := sqlquery.GetSqlInstance()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "fail",
			"message": "Failed to connect to db.",
		})
		return
	}
	defer db.Close()

	var latestVersion int
	err = db.QueryRow(sqlquery.SecretGetLatestVersionQuery, payload.SecretID).Scan(&latestVersion)
	if err != nil && err != sql.ErrNoRows {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "fail",
			"message": "Failed to fetch latest version",
		})
		return
	}

	newVersion := latestVersion + 1

	encryptedData, err := security.EncryptSecret(payload.SecretValue)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "Failed to encrypt secret.",
		})
		return
	}

	_, err = db.Exec(sqlquery.SecretVersionInsertQuery, payload.SecretID, newVersion, encryptedData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "Failed to insert new secret version",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": fmt.Sprintf("Added new version %d for secret with id %d", newVersion, payload.SecretID),
	})
}

func SecretGet(c *gin.Context) {
	secretId := c.Query("secretId")
	currentUser, err := user.GetCurrentUser(c)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "Failed to parse current user info.",
		})
		return
	}

	if secretId == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "Missing required parameter: secretId",
		})
		return
	}

	projectAccess, err := access.GetAccessForSecret(secretId, currentUser.Uid)

	if err != nil || !projectAccess.HasReadAccess {
		c.JSON(http.StatusForbidden, gin.H{
			"status":  "fail",
			"message": "Access denied to this resource",
		})
		return
	}

	db, err := sqlquery.GetSqlInstance()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "fail",
			"message": "Failed to connect to db.",
		})
		return
	}
	defer db.Close()

	rows, err := db.Query(sqlquery.SecretGetLatestValueQuery, secretId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "Failed to fetch secret",
		})
		return
	}
	defer rows.Close()

	var secretList []interfaces.SecretResponse

	for rows.Next() {
		var id, version int
		var name, value, environment string
		if err := rows.Scan(&id, &name, &version, &value, &environment); err != nil {
			log.Fatal(err)
		}

		decrytedValue, err := security.DecryptSecret(value)

		if err == nil {
			secretList = append(secretList, interfaces.SecretResponse{
				ID:          id,
				Name:        name,
				Version:     version,
				Value:       decrytedValue,
				Environment: environment,
			})
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   secretList,
	})
}

func SecretDelete(c *gin.Context) {
	var payload interfaces.SecretDeletePayload

	currentUser, err := user.GetCurrentUser(c)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "Failed to parse current user info.",
		})
		return
	}

	if err := c.BindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "Missing required parameter: id",
		})
		return
	}

	if payload.SecretID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "Missing required parameter: id",
		})
		return
	}

	projectAccess, err := access.GetAccessForSecret(payload.SecretID, currentUser.Uid)

	if err != nil || !projectAccess.HasDeleteAccess {
		c.JSON(http.StatusForbidden, gin.H{
			"status":  "fail",
			"message": "Access denied to this resource",
		})
		return
	}

	db, err := sqlquery.GetSqlInstance()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "fail",
			"message": "Failed to connect to db.",
		})
		return
	}
	defer db.Close()

	//Delete related entries first if you have a `secret_project_map` table
	_, err = db.Exec(sqlquery.DeleteSecretFromMapTable, payload.SecretID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "fail",
			"message": "Failed to delete secret-project mappings.",
		})
		return
	}

	// Delete versions from secret_version table
	_, err = db.Exec(sqlquery.DeleteSecretFromSecretVersionTable, payload.SecretID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "fail",
			"message": "Failed to delete secret.",
		})
		return
	}

	// Delete from secrets table
	result, err := db.Exec(sqlquery.DeleteSecretFromSecretTable, payload.SecretID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "fail",
			"message": "Failed to delete secret.",
		})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "fail",
			"message": "No secret found with the given ID.",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": fmt.Sprintf("Secret with ID %d deleted successfully.", payload.SecretID),
	})
}

func ListSecretByProjectUid(c *gin.Context) {
	projectUid := c.Query("projectuid")
	environment := c.Query("environment")

	if projectUid == "" || environment == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "Missing required parameter: projectuid and environment",
		})
		return
	}

	var secretList []interfaces.SecretResponse

	db, err := sqlquery.GetSqlInstance()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "Failed to connect to db.",
		})
		return
	}
	defer db.Close()

	rows, err := db.Query(sqlquery.GetAllSecretByProjectUidQuery, projectUid, environment)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "Failed to get secrets for project with project uid " + projectUid,
		})
		return
	}

	defer rows.Close()

	for rows.Next() {
		var id, version int
		var name, value, environment string
		if err := rows.Scan(&id, &name, &version, &value, &environment); err != nil {
			log.Fatal(err)
		}

		decrytedValue, err := security.DecryptSecret(value)
		if err == nil {
			secretList = append(secretList, interfaces.SecretResponse{
				ID:          id,
				Name:        name,
				Version:     version,
				Value:       decrytedValue,
				Environment: environment,
			})
		}
	}

	if err := rows.Err(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "Failed to get secrets for project with project uid " + projectUid,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   secretList,
	})

}

func ListSecretByUser(c *gin.Context) {
	currentUser, err := user.GetCurrentUser(c)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "Failed to parse current user info.",
		})
		return
	}

	var secretList []any

	db, err := sqlquery.GetSqlInstance()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "Failed to connect to db.",
		})
		return
	}
	defer db.Close()

	rows, err := db.Query(sqlquery.ListAllSecretByUser, currentUser.Uid)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "Failed to fetch secret list.",
		})
		return
	}

	for rows.Next() {
		var id int
		var name string
		if err := rows.Scan(&id, &name); err != nil {
			log.Fatal(err)
		}
		secretList = append(secretList, map[string]any{
			"id":   id,
			"name": name,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   secretList,
	})
}
