package api

import (
	"fmt"
	"log"
	"net/http"

	"cozeva.com/vault/interfaces"
	"cozeva.com/vault/pkg/access"
	"cozeva.com/vault/pkg/user"
	sqlquery "cozeva.com/vault/sql"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func WorkspaceCreate(c *gin.Context) {
	var newWorkspaceCreatePayload interfaces.WorkspaceCreatePayload

	currentUser, err := user.GetCurrentUser(c)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "Failed to parse current user info.",
		})
		return
	}

	if err := c.BindJSON(&newWorkspaceCreatePayload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "Missing required parameter: name",
		})
		return
	}

	if newWorkspaceCreatePayload.WorkspaceDisplayName == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "Property name must be a valid string.",
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

	tx, err := db.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "fail",
			"message": "Failed to start transaction",
		})
		return
	}

	result, err := tx.Exec(sqlquery.WorkspaceInsertQuery, newWorkspaceCreatePayload.WorkspaceDisplayName, uuid.NewString())

	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "Failed to create new workspace with name " + newWorkspaceCreatePayload.WorkspaceDisplayName,
		})
		return
	}

	workspaceId, _ := result.LastInsertId()

	_, err = tx.Exec(sqlquery.InsertUserWorkspaceAccessQuery, currentUser.Uid, workspaceId, 1, 1, 1)

	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "Failed to create new workspace with name " + newWorkspaceCreatePayload.WorkspaceDisplayName,
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
		"message": "Created new workspace with name " + newWorkspaceCreatePayload.WorkspaceDisplayName,
		"id":      workspaceId,
	})

}

func WorkspaceUpdate(c *gin.Context) {
	var newWorkspaceUpdatePayload interfaces.WorkspaceUpdatePayload

	currentUser, err := user.GetCurrentUser(c)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "Failed to parse current user info.",
		})
		return
	}

	if err := c.BindJSON(&newWorkspaceUpdatePayload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "Missing required parameter: name and id",
		})
		return
	}

	if newWorkspaceUpdatePayload.WorkspaceDisplayName == "" || newWorkspaceUpdatePayload.WorkspaceID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "Property name and id must be a valid string, integer respectively.",
		})
		return
	}

	workspaceAccess, err := access.GetAccessForWorkspace(newWorkspaceUpdatePayload.WorkspaceID, currentUser.Uid)

	if err != nil || !workspaceAccess.HasWriteAccess {
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

	_, err = db.Exec(sqlquery.WorkspaceUpdateQuery, newWorkspaceUpdatePayload.WorkspaceDisplayName, newWorkspaceUpdatePayload.WorkspaceID)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "Failed to update workspace with name " + newWorkspaceUpdatePayload.WorkspaceDisplayName,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Updated workspace with name " + newWorkspaceUpdatePayload.WorkspaceDisplayName,
	})

}

func WorkspaceGet(c *gin.Context) {
	workspaceId := c.Query("workspaceId")

	currentUser, err := user.GetCurrentUser(c)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "Failed to parse current user info.",
		})
		return
	}

	if workspaceId == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "Missing required parameter: workspaceId",
		})
		return
	}

	workspaceAccess, err := access.GetAccessForWorkspace(workspaceId, currentUser.Uid)

	if err != nil || !workspaceAccess.HasReadAccess {
		c.JSON(http.StatusForbidden, gin.H{
			"status":  "fail",
			"message": "Access denied to this resource",
		})
		return
	}

	var workspaceDetails []map[string]any

	db, err := sqlquery.GetSqlInstance()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "Failed to connect to db.",
		})
		return
	}
	defer db.Close()

	rows, err := db.Query(sqlquery.WorkspaceGetQuery, workspaceId)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "Failed to get workspace info for workspace with id " + workspaceId,
		})
		return
	}

	defer rows.Close()

	for rows.Next() {
		var id int
		var name string
		var uid string

		err = rows.Scan(&id, &name, &uid)
		if err != nil {
			log.Fatal(err)
		}

		workspaceDetails = append(workspaceDetails, map[string]any{
			"id":   id,
			"name": name,
			"uid":  uid,
		})
	}

	if err := rows.Err(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "Failed to get workspace info for workspace with id " + workspaceId,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   workspaceDetails,
	})

}

func WorkspaceDelete(c *gin.Context) {
	var payload interfaces.WorkspaceDeletePayload

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

	if payload.WorkspaceID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "Missing required parameter: id",
		})
		return
	}

	workspaceAccess, err := access.GetAccessForWorkspace(payload.WorkspaceID, currentUser.Uid)

	if err != nil || !workspaceAccess.HasDeleteAccess {
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

	//Delete related entries first if you have a `secret_workspace_map` table
	_, err = db.Exec(sqlquery.DeleteWorkspaceFromMapTable, payload.WorkspaceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "fail",
			"message": "Failed to delete secret-workspace mappings.",
		})
		return
	}

	// Delete from workspaces table
	result, err := db.Exec(sqlquery.DeleteWorkspaceFromWorkspaceTable, payload.WorkspaceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "fail",
			"message": "Failed to delete workspace.",
		})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "fail",
			"message": "No workspace found with the given ID.",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": fmt.Sprintf("Workspace with ID %d deleted successfully.", payload.WorkspaceID),
	})

}

func ListWorkspaceByUser(c *gin.Context) {
	currentUser, err := user.GetCurrentUser(c)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "Failed to parse current user info.",
		})
		return
	}

	var workspaceList []any

	db, err := sqlquery.GetSqlInstance()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "Failed to connect to db.",
		})
		return
	}
	defer db.Close()

	rows, err := db.Query(sqlquery.ListAllWorkspaceByUser, currentUser.Uid)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "Failed to fetch workspace list.",
		})
		return
	}

	for rows.Next() {
		var id int
		var name string
		if err := rows.Scan(&id, &name); err != nil {
			log.Fatal(err)
		}
		workspaceList = append(workspaceList, map[string]any{
			"id":   id,
			"name": name,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   workspaceList,
	})
}
