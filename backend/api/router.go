package api

import (
	"github.com/fernandoglatz/home-management/configs"
	"github.com/gin-gonic/gin"
)

func Setup() error {
	router := gin.Default()

	router.POST("/users", CreateUser)
	router.GET("/users/:id", GetUser)
	router.GET("/users", GetAllUsers)
	router.PUT("/users/:id", UpdateUser)
	router.DELETE("/users/:id", DeleteUser)

	listening := configs.ApplicationConfig.Server.Listening

	return router.Run(listening)
}
