package controller

import (
	"fernandoglatz/home-management/internal/core/entity"
	"fernandoglatz/home-management/internal/core/service"
	"sync"

	"github.com/gin-gonic/gin"
)

var eventController any
var eventControllerMutex sync.Mutex

type EventController[T entity.IEntity] struct {
	controller Controller[T]
}

func GetEventController[T entity.IEntity]() EventController[T] {
	eventControllerMutex.Lock()
	defer eventControllerMutex.Unlock()

	if eventController == nil {
		eventService := service.GetEventService[T]()

		eventController = EventController[T]{
			controller: GetController[T](&eventService),
		}
	}

	return eventController.(EventController[T])
}

// @Tags	event
// @Summary	Get events
// @Produce	json
// @Success	200	{array}		entity.Event
// @Failure	400	{object}	response.Response
// @Failure	500	{object}	response.Response
// @Router	/event [get]
func (eventController *EventController[T]) Get(ginCtx *gin.Context) {
	eventController.controller.Get(ginCtx)
}

// @Tags	event
// @Summary	Get event
// @Param	id		path	string  true "id"
// @Produce	json
// @Success	200	{object}	entity.Event
// @Failure	400	{object}	response.Response
// @Failure	500	{object}	response.Response
// @Router	/event/{id} [get]
func (eventController *EventController[T]) GetById(ginCtx *gin.Context) {
	eventController.controller.GetById(ginCtx)
}

// @Tags	event
// @Summary	Update event
// @Param	id		path	string  true "id"
// @Param	request	body	entity.Event true "body"
// @Accept	json
// @Produce	json
// @Success	200	{object}	entity.Event
// @Failure	400	{object}	response.Response
// @Failure	500	{object}	response.Response
// @Router		/event/{id} [post]
func (eventController *EventController[T]) Post(ginCtx *gin.Context) {
	eventController.controller.Post(ginCtx)
}

// @Tags	event
// @Summary	Create event
// @Param	request	body	entity.Event true "body"
// @Accept	json
// @Produce	json
// @Success	200	{object}	entity.Event
// @Failure	400	{object}	response.Response
// @Failure	500	{object}	response.Response
// @Router		/event [put]
func (eventController *EventController[T]) Put(ginCtx *gin.Context) {
	eventController.controller.Put(ginCtx)
}

// @Tags	event
// @Summary	Update event
// @Param	id		path	string  true "id"
// @Param	request	body	entity.Event true "body"
// @Accept	json
// @Produce	json
// @Success	200	{object}	entity.Event
// @Failure	400	{object}	response.Response
// @Failure	500	{object}	response.Response
// @Router		/event/{id} [put]
func (eventController *EventController[T]) PutById(ginCtx *gin.Context) {
	eventController.controller.PutById(ginCtx)
}

// @Tags	event
// @Summary	Update event
// @Param	id		path	string  true "id"
// @Param	request	body	entity.Event true "body"
// @Accept	json
// @Produce	json
// @Success	200	{object}	entity.Event
// @Failure	400	{object}	response.Response
// @Failure	500	{object}	response.Response
// @Router		/event/{id} [patch]
func (eventController *EventController[T]) Patch(ginCtx *gin.Context) {
	eventController.controller.Patch(ginCtx)
}

// @Tags	event
// @Summary	Delete event
// @Param	id		path	string  true "id"
// @Produce	json
// @Success	200
// @Failure	400	{object}	response.Response
// @Failure	500	{object}	response.Response
// @Router	/event/{id} [delete]
func (eventController *EventController[T]) DeleteById(ginCtx *gin.Context) {
	eventController.controller.DeleteById(ginCtx)
}

// @Tags		events
// @Summary	Check if event exists
// @Param		id	path	string	true	"Event ID"
// @Success	200
// @Failure	404
// @Failure	500
// @Router		/v1/events/{id} [head]
func (eventController *EventController[T]) Head(ginCtx *gin.Context) {
	eventController.controller.Head(ginCtx)
}
