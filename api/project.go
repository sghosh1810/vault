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

func ProjectCreate(c *gin.Context) {
	var newProjectsCreatePayload interfaces.ProjectsCreatePayload

	currentUser, err := user.GetCurrentUser(c)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "Failed to parse current user info.",
		})
		return
	}

	if err := c.BindJSON(&newProjectsCreatePayload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "Missing required parameter: name and workspace_id",
		})
		return
	}

	if newProjectsCreatePayload.ProjectDisplayName == "" {
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

	workspaceAccess, err := access.GetAccessForWorkspace(newProjectsCreatePayload.WorkspaceID, currentUser.Uid)

	if err != nil || !workspaceAccess.HasWriteAccess {
		c.JSON(http.StatusForbidden, gin.H{
			"status":  "fail",
			"message": "You don't have write access to the workspace associated with this project.",
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

	result, err := tx.Exec(sqlquery.ProjectInsertQuery, newProjectsCreatePayload.ProjectDisplayName, uuid.NewString(), newProjectsCreatePayload.WorkspaceID)

	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "Failed to create new project with name " + newProjectsCreatePayload.ProjectDisplayName,
		})
		return
	}

	projectId, _ := result.LastInsertId()

	_, err = tx.Exec(sqlquery.InsertWorkspaceProjectAccessQuery, newProjectsCreatePayload.WorkspaceID, projectId, 1, 1, 1)

	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "Failed to create new project with name " + newProjectsCreatePayload.ProjectDisplayName,
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
		"message": "Created new project with name " + newProjectsCreatePayload.ProjectDisplayName,
		"id":      projectId,
	})

}

func ProjectUpdate(c *gin.Context) {
	var newProjectsUpdatePayload interfaces.ProjectsUpdatePayload

	currentUser, err := user.GetCurrentUser(c)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "Failed to parse current user info.",
		})
		return
	}

	if err := c.BindJSON(&newProjectsUpdatePayload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "Missing required parameter: name,id and workspace_id",
		})
		return
	}

	if newProjectsUpdatePayload.ProjectDisplayName == "" || newProjectsUpdatePayload.ProjectID == 0 || newProjectsUpdatePayload.WorkspaceID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "Property name and id, workspace_id must be a valid string, integer and integer respectively.",
		})
		return
	}

	projectAccess, err := access.GetAccessForProject(newProjectsUpdatePayload.ProjectID, newProjectsUpdatePayload.WorkspaceID, currentUser.Uid)

	if err != nil || !projectAccess.HasWriteAccess {
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

	_, err = db.Exec(sqlquery.ProjectUpdateQuery, newProjectsUpdatePayload.ProjectDisplayName, newProjectsUpdatePayload.ProjectID, newProjectsUpdatePayload.WorkspaceID)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "Failed to update project with name " + newProjectsUpdatePayload.ProjectDisplayName,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Updated project with name " + newProjectsUpdatePayload.ProjectDisplayName,
	})

}

func ProjectGet(c *gin.Context) {
	projectId := c.Query("projectId")
	workspaceId := c.Query("workspaceId")

	currentUser, err := user.GetCurrentUser(c)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "Failed to parse current user info.",
		})
		return
	}

	if projectId == "" || workspaceId == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "Missing required parameter: projectId and workspaceId",
		})
		return
	}

	projectAccess, err := access.GetAccessForProject(projectId, workspaceId, currentUser.Uid)

	if err != nil || !projectAccess.HasReadAccess {
		c.JSON(http.StatusForbidden, gin.H{
			"status":  "fail",
			"message": "Access denied to this resource",
		})
		return
	}

	var projectDetails []map[string]any

	db, err := sqlquery.GetSqlInstance()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "Failed to connect to db.",
		})
		return
	}
	defer db.Close()

	rows, err := db.Query(sqlquery.ProjectGetQuery, projectId, workspaceId)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "Failed to get project info for project with id " + projectId,
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

		projectDetails = append(projectDetails, map[string]any{
			"id":   id,
			"name": name,
			"uid":  uid,
		})
	}

	if err := rows.Err(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "Failed to get project info for project with id " + projectId,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   projectDetails,
	})

}

func ProjectDelete(c *gin.Context) {
	var payload interfaces.ProjectsDeletePayload

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
			"message": "Missing required parameter: id and workspace_id",
		})
		return
	}

	if payload.ProjectID == 0 || payload.WorkspaceID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "Missing required parameter: id and workspace_id",
		})
		return
	}

	projectAccess, err := access.GetAccessForProject(payload.ProjectID, payload.WorkspaceID, currentUser.Uid)

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
	_, err = db.Exec(sqlquery.DeleteProjectFromMapTable, payload.ProjectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "fail",
			"message": "Failed to delete secret-project mappings.",
		})
		return
	}

	// Delete from projects table
	result, err := db.Exec(sqlquery.DeleteProjectFromProjectTable, payload.ProjectID, payload.WorkspaceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "fail",
			"message": "Failed to delete project.",
		})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "fail",
			"message": "No project found with the given ID.",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": fmt.Sprintf("Project with ID %d deleted successfully.", payload.ProjectID),
	})

}

func ListProjectByUser(c *gin.Context) {
	workspaceId := c.Query("workspaceId")
	currentUser, err := user.GetCurrentUser(c)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "Failed to parse current user info.",
		})
		return
	}

	var projectList []any

	db, err := sqlquery.GetSqlInstance()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "Failed to connect to db.",
		})
		return
	}
	defer db.Close()

	rows, err := db.Query(sqlquery.ListAllProjectByUser, currentUser.Uid, workspaceId)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "Failed to fetch project list.",
		})
		return
	}

	for rows.Next() {
		var id int
		var name string
		if err := rows.Scan(&id, &name); err != nil {
			log.Fatal(err)
		}
		projectList = append(projectList, map[string]any{
			"id":   id,
			"name": name,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   projectList,
	})
}
