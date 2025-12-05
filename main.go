package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"cozeva.com/vault/api"
	"cozeva.com/vault/config"
	"cozeva.com/vault/middleware"
	"cozeva.com/vault/pkg/crud"
	"cozeva.com/vault/pkg/eurekaclient"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if os.Getenv("GIN_MODE") != "release" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	port, _ := strconv.Atoi(config.GetConfigValue("api.config.port"))

	if os.Getenv("core.eureka.register") == "true" {

		ttl, _ := strconv.Atoi(os.Getenv("core.eureka.ttl"))
		isSsl, _ := strconv.ParseBool(os.Getenv("core.eureka.ssl"))

		eurekaclient.InitEurekaClient(
			[]string{
				fmt.Sprintf("%s/eureka", os.Getenv("core.eureka.eurekaHost")),
			},
			os.Getenv("api.config.hostname"),
			os.Getenv("core.eureka.serviceId"),
			os.Getenv("api.config.ip"),
			port,
			uint(ttl),
			isSsl,
		)
	}

	router := gin.Default()

	//get version of api from config or default to v1
	version := config.GetConfigValue("api.config.version")

	// winLogger := httplogger.NewLogger("http://localhost:4000/logs", "go-gin-service")

	// router.Use(middleware.WinstonLogger(winLogger))

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

	//route for user specific resource lists
	authGroup.GET("/user/secret/list", api.ListSecretByUser)
	authGroup.GET("/user/project/list", api.ListProjectByUser)

	//route  to map secret with a project
	authGroup.POST("/operation/map", api.MapSecretToProject)

	anonGroup := router.Group(apiBasePath, middleware.CheckIP)

	//route to get project specific data
	anonGroup.GET("/secret/list", api.ListSecretByProjectUid)

	router.Run(fmt.Sprintf(":%d", port))
}
