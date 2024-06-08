package controller

import (
	"fernandoglatz/home-management/internal/core/entity"
	"fernandoglatz/home-management/internal/core/service"
	"fernandoglatz/home-management/internal/infrastructure/repository"

	"github.com/gin-gonic/gin"
)

type EventController struct {
	controller Controller[*entity.Event]
}

func NewEventController() EventController {
	eventRepository := repository.NewRepository[*entity.Event](&entity.Event{})
	eventService := service.NewService[*entity.Event](*eventRepository)

	return EventController{
		controller: NewController[*entity.Event]("event", eventService),
	}
}

// @Tags	event
// @Summary	Get events
// @Produce	json
// @Success	200	{array}		entity.Event
// @Failure	400	{object}	response.Response
// @Failure	500	{object}	response.Response
// @Router	/event [get]
func (eventController *EventController) Get(ginCtx *gin.Context) {
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
func (eventController *EventController) GetById(ginCtx *gin.Context) {
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
func (eventController *EventController) Post(ginCtx *gin.Context) {
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
func (eventController *EventController) Put(ginCtx *gin.Context) {
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
func (eventController *EventController) PutById(ginCtx *gin.Context) {
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
func (eventController *EventController) Patch(ginCtx *gin.Context) {
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
func (eventController *EventController) DeleteById(ginCtx *gin.Context) {
	eventController.controller.DeleteById(ginCtx)
}

// @Tags		events
// @Summary	Check if event exists
// @Param		id	path	string	true	"Event ID"
// @Success	200
// @Failure	404
// @Failure	500
// @Router		/v1/events/{id} [head]
func (eventController *EventController) Head(ginCtx *gin.Context) {
	eventController.controller.Head(ginCtx)
}
