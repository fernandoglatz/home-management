package router

import (
	"context"
	_ "fernandoglatz/home-management/docs"
	"fernandoglatz/home-management/internal/controller"
	"fernandoglatz/home-management/internal/core/common/utils/constants"
	"fernandoglatz/home-management/internal/core/common/utils/log"
	"fernandoglatz/home-management/internal/infrastructure/config"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Setup(ctx context.Context, engine *gin.Engine) {
	log.Info(ctx).Msg("Configuring routes")

	contextPath := config.ApplicationConfig.Server.ContextPath
	router := engine.Group(contextPath)

	healthController := controller.NewHealthController()

	eventController := controller.NewEventController()
	routerEvent := router.Group("/event")
	routerEvent.GET(constants.EMPTY, eventController.Get)
	routerEvent.GET(":id", eventController.GetById)
	routerEvent.HEAD(":id", eventController.Head)
	routerEvent.PUT(constants.EMPTY, eventController.Put)
	routerEvent.PUT(":id", eventController.PutById)
	routerEvent.POST(":id", eventController.Post)
	routerEvent.PATCH(":id", eventController.Patch)
	routerEvent.DELETE(":id", eventController.DeleteById)

	router.GET("/health", healthController.Health)
	router.GET("/swagger-ui/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	log.Info(ctx).Msg("Routes configured")
}
