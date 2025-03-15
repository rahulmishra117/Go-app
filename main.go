package main

import (
	"github.com/gin-gonic/gin"
	"github.com/rahulmishra/go-crud-app/config"
	_ "github.com/rahulmishra/go-crud-app/docs"
	"github.com/rahulmishra/go-crud-app/models"
	"github.com/rahulmishra/go-crud-app/routes"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	config.ConnectDatabase()
	config.DB.AutoMigrate(&models.Item{})
	config.ConnectRedis()
	r := gin.Default()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	routes.SetupItemRoutes(r)
	r.Run(":9000")
}
