package api

import (
	"github.com/fernandoglatz/home-management/configs"
	"github.com/fernandoglatz/home-management/models"
	"github.com/gin-gonic/gin"
)

func Setup() error {
	router := gin.Default()

	userController := NewController[*models.User]()
	userPath := userController.service.BaseEntity.GetEntityName()

	router.POST("/api/"+userPath, userController.Create)
	router.GET("/api/"+userPath, userController.GetAll)
	router.GET("/api/"+userPath+"/:id", userController.Get)
	router.PUT("/api/"+userPath+"/:id", userController.Update)
	router.PATCH("/api/"+userPath+"/:id", userController.Patch)
	router.DELETE("/api/"+userPath+"/:id", userController.Delete)

	listening := configs.ApplicationConfig.Server.Listening

	return router.Run(listening)
}
