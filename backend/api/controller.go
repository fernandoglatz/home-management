package api

import (
	"net/http"
	"strings"

	"github.com/fernandoglatz/home-management/api/dto"
	"github.com/fernandoglatz/home-management/models"
	"github.com/fernandoglatz/home-management/services"
	"github.com/gin-gonic/gin"
)

type Controller[T models.IEntity] struct {
	service *services.Service[T]
}

func NewController[T models.IEntity]() *Controller[T] {
	controller := &Controller[T]{}
	controller.service = services.NewService[T]()
	return controller
}

func (controller *Controller[T]) Create(c *gin.Context) {
	var entity T
	if err := c.ShouldBindJSON(&entity); err != nil {
		createResponse(c, http.StatusBadRequest, "Bad JSON received", err)
		return
	}

	service := controller.service
	entityName := service.BaseEntity.GetEntityName()
	err := service.Save(entity)
	if err != nil {
		httpStatus := http.StatusInternalServerError

		if strings.Contains(err.Error(), "duplicate") {
			httpStatus = http.StatusConflict
		}

		createResponse(c, httpStatus, "Failed to create a new entry in "+entityName, err)
		return
	}

	id := entity.GetID()
	c.Header("Location", "/api/"+entityName+"/"+id)
	c.JSON(http.StatusCreated, entity)
}

func (controller *Controller[T]) Update(c *gin.Context) {
	var entity T
	id := c.Param("id")
	if err := c.ShouldBindJSON(&entity); err != nil {
		createResponse(c, http.StatusBadRequest, "Bad JSON received", err)
		return
	}

	service := controller.service
	entityName := service.BaseEntity.GetEntityName()

	oldEntity, err := service.FindByID(id)
	if err != nil {
		if strings.Contains(err.Error(), "no documents") {
			createResponse(c, http.StatusNotFound, "Entry to update not found in "+entityName, nil)
			return
		}

		createResponse(c, http.StatusInternalServerError, "Failed to update a entry in "+entityName, err)
		return
	}

	entity.SetID(oldEntity.GetID())
	entity.SetCreatedAt(oldEntity.GetCreatedAt())

	err = service.Save(entity)
	if err != nil {
		httpStatus := http.StatusInternalServerError

		if strings.Contains(err.Error(), "duplicate") {
			httpStatus = http.StatusConflict
		}

		createResponse(c, httpStatus, "Failed to update a entry in "+entityName, err)
		return
	}

	c.JSON(http.StatusOK, entity)
}

func (controller *Controller[T]) Patch(c *gin.Context) {
	var entity T
	id := c.Param("id")
	if err := c.ShouldBindJSON(&entity); err != nil {
		createResponse(c, http.StatusBadRequest, "Bad JSON received", err)
		return
	}

	service := controller.service
	entityName := service.BaseEntity.GetEntityName()

	entity, err := service.FindByID(id)
	if err != nil {
		if strings.Contains(err.Error(), "no documents") {
			createResponse(c, http.StatusNotFound, "Entry to patch not found in "+entityName, nil)
			return
		}

		createResponse(c, http.StatusInternalServerError, "Failed to patch a entry in "+entityName, err)
		return
	}

	createdAt := entity.GetCreatedAt()

	err = c.Bind(entity)
	if err != nil {
		createResponse(c, http.StatusInternalServerError, "Failed to patch a entry in "+entityName, err)
	}

	entity.SetCreatedAt(createdAt)

	err = service.Save(entity)
	if err != nil {
		httpStatus := http.StatusInternalServerError

		if strings.Contains(err.Error(), "duplicate") {
			httpStatus = http.StatusConflict
		}

		createResponse(c, httpStatus, "Failed to patch a entry in "+entityName, err)
		return
	}

	c.JSON(http.StatusOK, entity)
}

func (controller *Controller[T]) Delete(c *gin.Context) {
	id := c.Param("id")

	service := controller.service
	entityName := service.BaseEntity.GetEntityName()

	entity, err := service.FindByID(id)
	if err != nil {
		if strings.Contains(err.Error(), "no documents") {
			createResponse(c, http.StatusNotFound, "Entry to delete not found in "+entityName, nil)
			return
		}

		createResponse(c, http.StatusInternalServerError, "Failed to delete a entry in "+entityName, err)
		return
	}

	err = service.Delete(entity)
	if err != nil {
		httpStatus := http.StatusInternalServerError
		createResponse(c, httpStatus, "Failed to delete a entry in "+entityName, err)
		return
	}

	createResponse(c, http.StatusOK, "Deleted entry in "+entityName, nil)
}

func (controller *Controller[T]) Get(c *gin.Context) {
	id := c.Param("id")

	service := controller.service
	entityName := service.BaseEntity.GetEntityName()

	entity, err := service.FindByID(id)
	if err != nil {
		if strings.Contains(err.Error(), "no documents") {
			createResponse(c, http.StatusNotFound, "Entry not found in "+entityName, nil)
			return
		}

		createResponse(c, http.StatusInternalServerError, "Failed to find a entry in "+entityName, err)
		return
	}

	c.JSON(http.StatusOK, entity)
}

func (controller *Controller[T]) GetAll(c *gin.Context) {
	service := controller.service
	entityName := service.BaseEntity.GetEntityName()

	entity, err := service.FindAll()
	if err != nil {
		createResponse(c, http.StatusInternalServerError, "Failed to find entries in "+entityName, err)
		return
	}

	c.JSON(http.StatusOK, entity)
}

func createResponse(c *gin.Context, httpStatus int, message string, err error) {
	responseDTO := new(dto.ResponseDTO)
	responseDTO.Message = message

	if err != nil {
		responseDTO.Details = err.Error()
	}

	c.JSON(httpStatus, responseDTO)
}
