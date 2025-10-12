package crud

import "github.com/gin-gonic/gin"

func RegisterCRUDRoutes(group *gin.RouterGroup, resource string, handlers map[string]gin.HandlerFunc) {
	r := group.Group("/" + resource)
	r.POST("/create", handlers["create"])
	r.PATCH("/update", handlers["update"])
	r.DELETE("/delete", handlers["delete"])
	r.GET("/get", handlers["get"])
}
