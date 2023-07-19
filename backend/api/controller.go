package api

import (
	"encoding/json"
	"net/http"
	"reflect"
	"strings"

	"github.com/fernandoglatz/home-management/backend/api/dtos"
	_ "github.com/fernandoglatz/home-management/backend/docs"
	"github.com/fernandoglatz/home-management/backend/models"
	"github.com/fernandoglatz/home-management/backend/services"
	"github.com/fernandoglatz/home-management/backend/utils"
	"github.com/gin-gonic/gin"
)

type Controller[T models.IEntity, D any] struct {
	basePath string
	service  services.Service[T]
}

func NewController[T models.IEntity, D any](basePath string) Controller[T, D] {
	return Controller[T, D]{
		service:  services.NewService[T](),
		basePath: basePath,
	}
}

func (controller *Controller[T, D]) Create(c *gin.Context) {
	var dto D
	var entity T
	json.Unmarshal([]byte("{}"), &entity) //new instance from generic

	err := c.ShouldBindJSON(&dto)
	if err != nil {
		createResponse(c, http.StatusBadRequest, "Bad JSON received", err)
		return
	}

	service := controller.service
	basePath := controller.basePath

	err = utils.CopyStructFields(dto, entity)
	if err != nil {
		createResponse(c, http.StatusInternalServerError, "Failed to create a new entry in "+basePath, err)
		return
	}

	err = service.Save(entity)
	if err != nil {
		httpStatus := http.StatusInternalServerError

		if strings.Contains(err.Error(), "duplicate") {
			httpStatus = http.StatusConflict
		}

		createResponse(c, httpStatus, "Failed to create a new entry in "+basePath, err)
		return
	}

	id := entity.GetID()
	c.Header("Location", "/api/v1/"+basePath+"/"+id)
	c.JSON(http.StatusCreated, entity)
}

func (controller *Controller[T, D]) Update(c *gin.Context) {
	var dto D
	var entity T
	json.Unmarshal([]byte("{}"), &entity) //new instance from generic

	id := c.Param("id")
	err := c.ShouldBindJSON(&dto)
	if err != nil {
		createResponse(c, http.StatusBadRequest, "Bad JSON received", err)
		return
	}

	service := controller.service
	basePath := controller.basePath

	err = utils.CopyStructFields(dto, entity)
	if err != nil {
		createResponse(c, http.StatusInternalServerError, "Failed to create a new entry in "+basePath, err)
		return
	}

	oldEntity, err := service.FindByID(id)
	if err != nil {
		if strings.Contains(err.Error(), "no documents") {
			createResponse(c, http.StatusNotFound, "Entry to update not found in "+basePath, nil)
			return
		}

		createResponse(c, http.StatusInternalServerError, "Failed to update a entry in "+basePath, err)
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

		createResponse(c, httpStatus, "Failed to update a entry in "+basePath, err)
		return
	}

	c.JSON(http.StatusOK, entity)
}

func (controller *Controller[T, D]) Patch(c *gin.Context) {
	id := c.Param("id")

	var jsonData map[string]interface{}
	err := c.ShouldBindJSON(&jsonData)
	if err != nil {
		createResponse(c, http.StatusBadRequest, "Bad JSON received", err)
	}

	service := controller.service
	basePath := controller.basePath

	entity, err := service.FindByID(id)
	if err != nil {
		if strings.Contains(err.Error(), "no documents") {
			createResponse(c, http.StatusNotFound, "Entry to patch not found in "+basePath, nil)
			return
		}

		createResponse(c, http.StatusInternalServerError, "Failed to patch a entry in "+basePath, err)
		return
	}

	entityValue := reflect.ValueOf(entity).Elem()
	entityType := reflect.TypeOf(entity).Elem()
	for i := 0; i < entityType.NumField(); i++ {
		fieldType := entityType.Field(i)
		fieldValue := entityValue.Field(i)
		fieldJsonTag := fieldType.Tag.Get("json")
		jsonTag := strings.Replace(fieldJsonTag, ",omitempty", "", -1)

		value := jsonData[jsonTag]
		if value != nil {
			fieldValue.Set(reflect.ValueOf(value))
		}
	}

	err = service.Save(entity)
	if err != nil {
		httpStatus := http.StatusInternalServerError

		if strings.Contains(err.Error(), "duplicate") {
			httpStatus = http.StatusConflict
		}

		createResponse(c, httpStatus, "Failed to patch a entry in "+basePath, err)
		return
	}

	c.JSON(http.StatusOK, entity)
}

func (controller *Controller[T, D]) Delete(c *gin.Context) {
	id := c.Param("id")

	service := controller.service
	basePath := controller.basePath

	entity, err := service.FindByID(id)
	if err != nil {
		if strings.Contains(err.Error(), "no documents") {
			createResponse(c, http.StatusNotFound, "Entry to delete not found in "+basePath, nil)
			return
		}

		createResponse(c, http.StatusInternalServerError, "Failed to delete a entry in "+basePath, err)
		return
	}

	err = service.Delete(entity)
	if err != nil {
		httpStatus := http.StatusInternalServerError
		createResponse(c, httpStatus, "Failed to delete a entry in "+basePath, err)
		return
	}

	createResponse(c, http.StatusOK, "Deleted entry in "+basePath, nil)
}

func (controller *Controller[T, D]) Head(c *gin.Context) {
	id := c.Param("id")

	service := controller.service
	basePath := controller.basePath

	_, err := service.FindByID(id)
	if err != nil {
		if strings.Contains(err.Error(), "no documents") {
			createResponse(c, http.StatusNotFound, "Entry not found in "+basePath, nil)
			return
		}

		createResponse(c, http.StatusInternalServerError, "Failed to find a entry in "+basePath, err)
		return
	}

	c.Status(http.StatusOK)
}

func (controller *Controller[T, D]) Get(c *gin.Context) {
	id := c.Param("id")

	service := controller.service
	basePath := controller.basePath

	entity, err := service.FindByID(id)
	if err != nil {
		if strings.Contains(err.Error(), "no documents") {
			createResponse(c, http.StatusNotFound, "Entry not found in "+basePath, nil)
			return
		}

		createResponse(c, http.StatusInternalServerError, "Failed to find a entry in "+basePath, err)
		return
	}

	c.JSON(http.StatusOK, entity)
}

func (controller *Controller[T, D]) GetAll(c *gin.Context) {
	service := controller.service
	basePath := controller.basePath

	entities, err := service.FindAll()
	if err != nil {
		createResponse(c, http.StatusInternalServerError, "Failed to find entries in "+basePath, err)
		return
	}

	c.JSON(http.StatusOK, entities)
}

func (controller *Controller[T, D]) Options(c *gin.Context) {
	c.Header("Allow", "GET, POST, PUT, PATCH,  DELETE, HEAD, OPTIONS")
	c.Status(http.StatusOK)
}

func createResponse(c *gin.Context, httpStatus int, message string, err error) {
	responseDTO := dtos.ResponseDTO{
		Message: message,
	}

	if err != nil {
		responseDTO.Details = err.Error()
	}

	c.JSON(httpStatus, responseDTO)
}
