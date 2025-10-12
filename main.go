package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"shounak.me/configmanager/api"
	"shounak.me/configmanager/middleware"
	"shounak.me/configmanager/pkg/crud"
)

func main() {
	if os.Getenv("GIN_MODE") != "release" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}
	router := gin.Default()

	//get version of api from config or default to v1
	version := os.Getenv("api.version")

	if version == "" {
		version = "v1"
	}

	apiBasePath := "api/" + version

	// router.Use(middlewares.Logger())

	authGroup := router.Group(apiBasePath, middleware.CheckIP, middleware.CheckAuth)

	crud.RegisterCRUDRoutes(authGroup, "project", map[string]gin.HandlerFunc{
		"create": api.ProjectCreate,
		"update": api.ProjectUpdate,
		"delete": api.ProjectDelete,
		"get":    api.ProjectGet,
	})

	crud.RegisterCRUDRoutes(authGroup, "secret", map[string]gin.HandlerFunc{
		"create": api.SecretCreate,
		"update": api.SecretUpdate,
		"delete": api.SecretDelete,
		"get":    api.SecretGet,
	})

	//route  to map secret with a project
	authGroup.POST("/operation/map", api.MapSecretToProject)

	anonGroup := router.Group(apiBasePath, middleware.CheckIP)

	//route to get project specific data
	anonGroup.GET("/secret/list", api.GetAllSecretByProjectUid)

	router.Run()
}
