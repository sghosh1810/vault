package crud

import "github.com/gin-gonic/gin"

func RegisterCRUDRoutes(group *gin.RouterGroup, resource string, handlers map[string]gin.HandlerFunc) {
	r := group.Group("/" + resource)
	if handlers["create"] != nil {
		r.POST("/create", handlers["create"])
	}
	if handlers["update"] != nil {
		r.PATCH("/update", handlers["update"])
	}
	if handlers["delete"] != nil {
		r.DELETE("/delete", handlers["delete"])
	}
	if handlers["get"] != nil {
		r.GET("/get", handlers["get"])
	}
	if handlers["list"] != nil {
		r.GET("/list", handlers["list"])
	}
}
