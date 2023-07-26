package api

import (
	"github.com/fernandoglatz/home-management/backend/configs"
	_ "github.com/fernandoglatz/home-management/backend/docs"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Setup() error {
	router := gin.Default()
	router.Use(CORSMiddleware())

	routerV1 := router.Group("/api/v1")

	userController := NewUserController()
	routerV1User := routerV1.Group("/users")
	routerV1User.OPTIONS("", userController.Options)
	routerV1User.POST("", userController.Create)
	routerV1User.GET("", userController.GetAll)
	routerV1User.GET("/:id", userController.Get)
	routerV1User.HEAD("/:id", userController.Head)
	routerV1User.PUT("/:id", userController.Update)
	routerV1User.PATCH("/:id", userController.Patch)
	routerV1User.DELETE("/:id", userController.Delete)

	homeController := NewHomeController()
	routerV1Home := routerV1.Group("/homes")
	routerV1Home.OPTIONS("", homeController.Options)
	routerV1Home.POST("", homeController.Create)
	routerV1Home.GET("", homeController.GetAll)
	routerV1Home.GET("/:id", homeController.Get)
	routerV1Home.HEAD("/:id", homeController.Head)
	routerV1Home.PUT("/:id", homeController.Update)
	routerV1Home.PATCH("/:id", homeController.Patch)
	routerV1Home.DELETE("/:id", homeController.Delete)

	deviceController := NewDeviceController()
	routerV1Device := routerV1.Group("/devices")
	routerV1Device.OPTIONS("", deviceController.Options)
	routerV1Device.POST("", deviceController.Create)
	routerV1Device.GET("", deviceController.GetAll)
	routerV1Device.GET("/:id", deviceController.Get)
	routerV1Device.HEAD("/:id", deviceController.Head)
	routerV1Device.PUT("/:id", deviceController.Update)
	routerV1Device.PATCH("/:id", deviceController.Patch)
	routerV1Device.DELETE("/:id", deviceController.Delete)

	eventController := NewEventController()
	routerV1Event := routerV1.Group("/events")
	routerV1Event.OPTIONS("", eventController.Options)
	routerV1Event.POST("", eventController.Create)
	routerV1Event.GET("", eventController.GetAll)
	routerV1Event.GET("/:id", eventController.Get)
	routerV1Event.HEAD("/:id", eventController.Head)
	routerV1Event.PUT("/:id", eventController.Update)
	routerV1Event.PATCH("/:id", eventController.Patch)
	routerV1Event.DELETE("/:id", eventController.Delete)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	listening := configs.ApplicationConfig.Server.Listening
	return router.Run(listening)
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		allowOrigin := configs.ApplicationConfig.Server.AllowOrigin
		c.Header("Access-Control-Allow-Origin", allowOrigin)
		c.Header("Access-Control-Allow-Headers", "*")
		c.Header("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, HEAD, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
