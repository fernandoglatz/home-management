package controller

import (
	"context"
	"encoding/json"
	"net/http"
	"reflect"
	"strconv"
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

type IController[T entity.IEntity, RQ any] interface {
	Get(ginCtx *gin.Context)
	GetById(ginCtx *gin.Context)
	Post(ginCtx *gin.Context)
	Put(ginCtx *gin.Context)
	PutById(ginCtx *gin.Context)
	DeleteById(ginCtx *gin.Context)
	Patch(ginCtx *gin.Context)
	Head(ginCtx *gin.Context)
}

const HEADER_PAGE = "page"
const HEADER_LIMIT = "limit"

var controllers map[string]any
var controllerMutex sync.Mutex

type Controller[T entity.IEntity, RQ any] struct {
	service service_port.IService[T]
}

func GetGenericController[T entity.IEntity, RQ any]() Controller[T, RQ] {
	entity := utils.Instance[T]()
	typeName := utils.GetTypeName(entity)
	controller := controllers[typeName]

	if controller == nil {
		service := service.GetGenericService[T]()
		controller = GetController[T, RQ](&service)
	}

	return controller.(Controller[T, RQ])
}

func GetController[T entity.IEntity, RQ any](service service_port.IService[T]) Controller[T, RQ] {
	entity := utils.Instance[T]()
	typeName := utils.GetTypeName(entity)

	controllerMutex.Lock()
	defer controllerMutex.Unlock()

	if controllers == nil {
		controllers = make(map[string]any)
	}

	controller := controllers[typeName]

	if controller == nil {
		controller = Controller[T, RQ]{
			service: service,
		}

		controllers[typeName] = controller
	}

	return controller.(Controller[T, RQ])
}

func (controller Controller[T, RQ]) Get(ginCtx *gin.Context) {
	ctx := GetContext(ginCtx)
	pageStr, errw := GetQuery(ginCtx, HEADER_PAGE, true)
	if errw != nil {
		HandleError(ctx, ginCtx, errw)
		return
	}

	limitStr, errw := GetQuery(ginCtx, HEADER_LIMIT, true)
	if errw != nil {
		HandleError(ctx, ginCtx, errw)
		return
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		HandleError(ctx, ginCtx, &exceptions.WrappedError{
			BaseError: exceptions.InvalidParameterFormat,
			Error:     err,
		})
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		HandleError(ctx, ginCtx, &exceptions.WrappedError{
			BaseError: exceptions.InvalidParameterFormat,
			Error:     err,
		})
		return
	}

	devices, errw := controller.service.GetAll(ctx, page, limit)
	if errw != nil {
		HandleError(ctx, ginCtx, errw)
		return
	}

	ginCtx.JSON(http.StatusOK, devices)
}

func (controller Controller[T, RQ]) GetById(ginCtx *gin.Context) {
	ctx := GetContext(ginCtx)
	id, err := GetParameter(ginCtx, constants.ID, true)
	if err != nil {
		HandleError(ctx, ginCtx, err)
		return
	}

	device, err := controller.service.Get(ctx, id)
	if err != nil {
		HandleError(ctx, ginCtx, err)
		return
	}

	ginCtx.JSON(http.StatusOK, device)
}

func (controller Controller[T, RQ]) Post(ginCtx *gin.Context) {
	ctx := GetContext(ginCtx)
	id, err := GetParameter(ginCtx, constants.ID, true)
	if err != nil {
		HandleError(ctx, ginCtx, err)
		return
	}

	controller.save(ginCtx, &id, false)
}

func (controller Controller[T, RQ]) Put(ginCtx *gin.Context) {
	controller.save(ginCtx, nil, true)
}

func (controller Controller[T, RQ]) PutById(ginCtx *gin.Context) {
	ctx := GetContext(ginCtx)
	id, err := GetParameter(ginCtx, constants.ID, true)
	if err != nil {
		HandleError(ctx, ginCtx, err)
		return
	}

	controller.save(ginCtx, &id, true)
}

func (controller Controller[T, RQ]) DeleteById(ginCtx *gin.Context) {
	ctx := GetContext(ginCtx)
	id, err := GetParameter(ginCtx, constants.ID, true)
	if err != nil {
		HandleError(ctx, ginCtx, err)
		return
	}

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

func (controller Controller[T, RQ]) save(ginCtx *gin.Context, id *string, override bool) {
	ctx := GetContext(ginCtx)

	entity := utils.Instance[T]()
	request := utils.Instance[RQ]()

	var errw *exceptions.WrappedError

	err := ginCtx.ShouldBindJSON(&request)
	if err != nil {
		HandleError(ctx, ginCtx, &exceptions.WrappedError{
			BaseError: exceptions.InvalidJSON,
			Error:     err,
		})
		return
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		HandleError(ctx, ginCtx, &exceptions.WrappedError{
			BaseError: exceptions.InvalidJSON,
			Error:     err,
		})
		return
	}

	err = json.Unmarshal(jsonData, &entity)
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

	errw = controller.service.Save(ctx, entity)
	if errw != nil {
		HandleError(ctx, ginCtx, errw)
		return

	} else {
		ginCtx.JSON(http.StatusOK, entity)
	}
}

func (controller Controller[T, RQ]) Patch(ginCtx *gin.Context) {
	ctx := GetContext(ginCtx)
	id, errw := GetParameter(ginCtx, constants.ID, true)
	if errw != nil {
		HandleError(ctx, ginCtx, errw)
		return
	}

	var jsonData map[string]any
	err := ginCtx.ShouldBindJSON(&jsonData)
	if err != nil {
		errw := &exceptions.WrappedError{
			BaseError: exceptions.InvalidJSON,
			Error:     err,
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
		jsonTag := strings.Replace(fieldJsonTag, ",omitempty", constants.EMPTY, constants.MINUS_ONE)

		value := jsonData[jsonTag]
		if value != nil {
			fieldValue.Set(reflect.ValueOf(value))
		}
	}

	errw = controller.service.Save(ctx, entity)
	if errw != nil {
		HandleError(ctx, ginCtx, errw)
		return

	} else {
		ginCtx.JSON(http.StatusOK, entity)
	}
}

func (controller Controller[T, RQ]) Head(ginCtx *gin.Context) {
	ctx := GetContext(ginCtx)
	id, err := GetParameter(ginCtx, constants.ID, true)
	if err != nil {
		HandleError(ctx, ginCtx, err)
		return
	}

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

func GetHeader(ginCtx *gin.Context, name string, required bool) (string, *exceptions.WrappedError) {
	header := ginCtx.Request.Header.Get(name)

	if utils.IsEmptyStr(header) && required {
		return header, &exceptions.WrappedError{
			Message:   "Header [" + name + "] not found",
			BaseError: exceptions.HeaderNotFound,
		}
	}

	return header, nil
}

func GetParameter(ginCtx *gin.Context, name string, required bool) (string, *exceptions.WrappedError) {
	parameter := ginCtx.Param(name)

	if utils.IsEmptyStr(parameter) && required {
		return parameter, &exceptions.WrappedError{
			Message:   "Parameter [" + name + "] not found",
			BaseError: exceptions.ParameterNotFound,
		}
	}

	return parameter, nil
}

func GetQuery(ginCtx *gin.Context, name string, required bool) (string, *exceptions.WrappedError) {
	query := ginCtx.Query(name)

	if utils.IsEmptyStr(query) && required {
		return query, &exceptions.WrappedError{
			Message:   "Query [" + name + "] not found",
			BaseError: exceptions.QueryNotFound,
		}
	}

	return query, nil
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
