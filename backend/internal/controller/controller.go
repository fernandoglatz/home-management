package controller

import (
	"context"
	"encoding/json"
	"net/http"
	"reflect"
	"strings"
	"sync"

	_ "fernandoglatz/home-management/docs"
	"fernandoglatz/home-management/internal/core/common/utils"
	"fernandoglatz/home-management/internal/core/common/utils/constants"
	"fernandoglatz/home-management/internal/core/common/utils/exceptions"
	"fernandoglatz/home-management/internal/core/common/utils/log"
	"fernandoglatz/home-management/internal/core/entity"
	"fernandoglatz/home-management/internal/core/model/response"
	service_port "fernandoglatz/home-management/internal/core/port/service"
	"fernandoglatz/home-management/internal/core/service"

	"github.com/gin-gonic/gin"
)

type IController[T entity.IEntity] interface {
	Get(ginCtx *gin.Context)
	GetById(ginCtx *gin.Context)
	Post(ginCtx *gin.Context)
	Put(ginCtx *gin.Context)
	PutById(ginCtx *gin.Context)
	DeleteById(ginCtx *gin.Context)
	Patch(ginCtx *gin.Context)
	Head(ginCtx *gin.Context)
}

var controllers map[string]any
var controllerMutex sync.Mutex

type Controller[T entity.IEntity] struct {
	service service_port.IService[T]
}

func GetGenericController[T entity.IEntity]() Controller[T] {
	var entity T
	typeName := utils.GetTypeName(entity)
	controller := controllers[typeName]

	if controller == nil {
		service := service.GetGenericService[T]()
		controller = GetController[T](&service)
	}

	return controller.(Controller[T])
}

func GetController[T entity.IEntity](service service_port.IService[T]) Controller[T] {
	var entity T
	typeName := utils.GetTypeName(entity)

	controllerMutex.Lock()
	defer controllerMutex.Unlock()

	if controllers == nil {
		controllers = make(map[string]any)
	}

	controller := controllers[typeName]

	if controller == nil {
		controller = Controller[T]{
			service: service,
		}

		controllers[typeName] = controller
	}

	return controller.(Controller[T])
}

func (controller *Controller[T]) Get(ginCtx *gin.Context) {
	ctx := GetContext(ginCtx)

	devices, err := controller.service.GetAll(ctx)
	if err != nil {
		HandleError(ctx, ginCtx, err)
		return
	}

	ginCtx.JSON(http.StatusOK, devices)
}

func (controller *Controller[T]) GetById(ginCtx *gin.Context) {
	ctx := GetContext(ginCtx)
	id := ginCtx.Param("id")

	device, err := controller.service.Get(ctx, id)
	if err != nil {
		HandleError(ctx, ginCtx, err)
		return
	}

	ginCtx.JSON(http.StatusOK, device)
}

func (controller *Controller[T]) Post(ginCtx *gin.Context) {
	id := ginCtx.Param(constants.ID)
	controller.save(ginCtx, &id, false)
}

func (controller *Controller[T]) Put(ginCtx *gin.Context) {
	controller.save(ginCtx, nil, true)
}

func (controller *Controller[T]) PutById(ginCtx *gin.Context) {
	id := ginCtx.Param(constants.ID)
	controller.save(ginCtx, &id, true)
}

func (controller *Controller[T]) DeleteById(ginCtx *gin.Context) {
	ctx := GetContext(ginCtx)
	id := ginCtx.Param("id")

	device, err := controller.service.Get(ctx, id)
	if err != nil {
		HandleError(ctx, ginCtx, err)
		return
	}

	err = controller.service.Remove(ctx, device)
	if err != nil {
		HandleError(ctx, ginCtx, err)
	} else {
		ginCtx.Status(http.StatusNoContent)
	}
}

func (controller *Controller[T]) save(ginCtx *gin.Context, id *string, override bool) {
	ctx := GetContext(ginCtx)

	var entity T
	json.Unmarshal([]byte("{}"), &entity) //new instance from generic

	var errw *exceptions.WrappedError

	err := ginCtx.ShouldBindJSON(&entity)
	if err != nil {
		HandleError(ctx, ginCtx, &exceptions.WrappedError{
			BaseError: exceptions.InvalidJSON,
			Error:     err,
		})
		return
	}

	if id != nil {
		entity, errw = controller.service.Get(ctx, *id)
		if errw != nil && !override {
			HandleError(ctx, ginCtx, errw)
			return
		}
		entity.SetID(*id)
	}

	errw = controller.service.Save(ctx, &entity)
	if errw != nil {
		HandleError(ctx, ginCtx, errw)
		return

	} else {
		ginCtx.JSON(http.StatusOK, entity)
	}
}

func (controller *Controller[T]) Patch(ginCtx *gin.Context) {
	ctx := GetContext(ginCtx)
	id := ginCtx.Param("id")

	var jsonData map[string]any
	err := ginCtx.ShouldBindJSON(&jsonData)
	if err != nil {
		errw := &exceptions.WrappedError{
			Error: err,
		}

		HandleError(ctx, ginCtx, errw)
		return
	}

	service := controller.service

	entity, errw := service.Get(ctx, id)
	if errw != nil {
		HandleError(ctx, ginCtx, errw)
		return
	}

	entityValue := reflect.ValueOf(entity).Elem()
	entityType := reflect.TypeOf(entity).Elem()
	for i := 0; i < entityType.NumField(); i++ {
		fieldType := entityType.Field(i)
		fieldValue := entityValue.Field(i)
		fieldJsonTag := fieldType.Tag.Get("json")
		jsonTag := strings.Replace(fieldJsonTag, ",omitempty", constants.EMPTY, -1)

		value := jsonData[jsonTag]
		if value != nil {
			fieldValue.Set(reflect.ValueOf(value))
		}
	}

	errw = controller.service.Save(ctx, &entity)
	if errw != nil {
		HandleError(ctx, ginCtx, errw)
		return

	} else {
		ginCtx.JSON(http.StatusOK, entity)
	}
}

func (controller *Controller[T]) Head(ginCtx *gin.Context) {
	ctx := GetContext(ginCtx)
	id := ginCtx.Param("id")

	service := controller.service

	_, errw := service.Get(ctx, id)
	if errw != nil {
		HandleError(ctx, ginCtx, errw)
		return
	}

	ginCtx.Status(http.StatusOK)
}

func GetContext(ginCtx *gin.Context) context.Context {
	return ginCtx.Request.Context()
}

func GetHeader(ctx *gin.Context, name string, required bool) (string, *exceptions.WrappedError) {
	header := ctx.Request.Header.Get(name)

	if utils.IsEmptyStr(header) && required {
		return header, &exceptions.WrappedError{
			Message:   "Header [" + name + "] não encontrado",
			BaseError: exceptions.HeaderNotFound,
		}
	}

	return header, nil
}

func HandleError(ctx context.Context, ginCtx *gin.Context, err *exceptions.WrappedError) {
	request := ginCtx.Request
	method := request.Method
	path := request.URL.Path

	code := err.GetCode()
	message := err.GetMessage()
	httpStatus := http.StatusBadRequest

	if err.Error != nil && len(err.BaseError.Code) == constants.ZERO {
		httpStatus = http.StatusInternalServerError
		code = http.StatusText(httpStatus)
		log.Error(ctx).Msg("[" + method + "] " + path + " - " + code + " - " + message)

	} else {
		log.Warn(ctx).Msg("[" + method + "] " + path + " - " + code + " - " + message)
	}

	if err.BaseError == exceptions.RecordNotFound {
		httpStatus = http.StatusNotFound
	}

	if err.BaseError == exceptions.DuplicatedRecord {
		httpStatus = http.StatusConflict
	}

	ginCtx.JSON(httpStatus, response.Response{
		Code:    code,
		Message: message,
	})
}
