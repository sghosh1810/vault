package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"shounak.me/configmanager/interfaces"
	sqlquery "shounak.me/configmanager/sql"
)

func ProjectCreate(c *gin.Context) {
	var newProjectsCreatePayload interfaces.ProjectsCreatePayload

	if err := c.BindJSON(&newProjectsCreatePayload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "Failed to decode json.",
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
	defer db.Close()

	_, err = db.Exec(sqlquery.ProjectInsertQuery, newProjectsCreatePayload.ProjectDisplayName, uuid.NewString())

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "Failed to create new project with name " + newProjectsCreatePayload.ProjectDisplayName,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Created new project with name " + newProjectsCreatePayload.ProjectDisplayName,
	})

}

func ProjectUpdate(c *gin.Context) {
	var newProjectsUpdatePayload interfaces.ProjectsUpdatePayload

	if err := c.BindJSON(&newProjectsUpdatePayload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "Failed to decode json.",
		})
		return
	}

	if newProjectsUpdatePayload.ProjectDisplayName == "" || newProjectsUpdatePayload.ProjectID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "Property name and id must be a valid string, integer respectively.",
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

	_, err = db.Exec(sqlquery.ProjectUpdateQuery, newProjectsUpdatePayload.ProjectDisplayName, newProjectsUpdatePayload.ProjectID)

	if err != nil {
		fmt.Println(err)
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

	rows, err := db.Query(sqlquery.ProjectGetQuery, projectId)

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

	if err := c.BindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "Failed to decode JSON",
		})
		return
	}

	if payload.ProjectID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "Missing required parameter: id",
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
	result, err := db.Exec(sqlquery.DeleteProjectFromProjectTable, payload.ProjectID)
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
