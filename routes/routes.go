package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/rahulmishra/go-crud-app/controllers"
)

func SetupItemRoutes(router *gin.Engine) {
	itemRoutes := router.Group("/items")
	{
		itemRoutes.POST("/", controllers.CreateItem)
		itemRoutes.GET("/", controllers.GetAllItems)
		itemRoutes.GET("/:id", controllers.GetItemByID)
		itemRoutes.PUT("/:id", controllers.UpdateItem)
		itemRoutes.DELETE("/:id", controllers.DeleteItem)
	}
}
