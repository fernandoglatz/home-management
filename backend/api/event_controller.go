package api

import (
	"github.com/fernandoglatz/home-management/backend/api/dtos"
	"github.com/fernandoglatz/home-management/backend/models"
	"github.com/gin-gonic/gin"
)

type EventController struct {
	controller Controller[*models.Event, *dtos.EventDTO]
}

func NewEventController() EventController {
	return EventController{
		controller: NewController[*models.Event, *dtos.EventDTO]("events"),
	}
}

//	@Tags		events
//	@Summary	Create event
//	@Param		request	body	dtos.EventDTO	true	"body"
//	@Accept		json
//	@Produce	json
//	@Success	200	{object}	models.Event
//	@Failure	400	{object}	dtos.ResponseDTO
//	@Failure	500	{object}	dtos.ResponseDTO
//	@Router		/v1/events [post]
func (eventController *EventController) Create(c *gin.Context) {
	eventController.controller.Create(c)
}

//	@Tags		events
//	@Summary	Update event
//	@Param		id		path	string			true	"Event ID"
//	@Param		request	body	dtos.EventDTO	true	"body"
//	@Accept		json
//	@Produce	json
//	@Success	200	{object}	models.Event
//	@Failure	400	{object}	dtos.ResponseDTO
//	@Failure	500	{object}	dtos.ResponseDTO
//	@Router		/v1/events/{id} [put]
func (eventController *EventController) Update(c *gin.Context) {
	eventController.controller.Update(c)
}

//	@Tags		events
//	@Summary	Partial update event
//	@Param		id		path	string			true	"Event ID"
//	@Param		request	body	dtos.EventDTO	true	"body"
//	@Accept		json
//	@Produce	json
//	@Success	200	{object}	models.Event
//	@Failure	400	{object}	dtos.ResponseDTO
//	@Failure	500	{object}	dtos.ResponseDTO
//	@Router		/v1/events/{id} [patch]
func (eventController *EventController) Patch(c *gin.Context) {
	eventController.controller.Patch(c)
}

//	@Tags		events
//	@Summary	Delete event
//	@Param		id	path	string	true	"Event ID"
//	@Produce	json
//	@Failure	200	{object}	dtos.ResponseDTO
//	@Failure	400	{object}	dtos.ResponseDTO
//	@Failure	500	{object}	dtos.ResponseDTO
//	@Router		/v1/events/{id} [delete]
func (eventController *EventController) Delete(c *gin.Context) {
	eventController.controller.Delete(c)
}

//	@Tags		events
//	@Summary	Check if event exists
//	@Param		id	path	string	true	"Event ID"
//	@Router		/v1/events/{id} [head]
func (eventController *EventController) Head(c *gin.Context) {
	eventController.controller.Head(c)
}

//	@Tags		events
//	@Summary	Find event
//	@Param		id	path	string	true	"Event ID"
//	@Produce	json
//	@Success	200	{object}	models.Event
//	@Failure	400	{object}	dtos.ResponseDTO
//	@Failure	500	{object}	dtos.ResponseDTO
//	@Router		/v1/events/{id} [get]
func (eventController *EventController) Get(c *gin.Context) {
	eventController.controller.Get(c)
}

//	@Tags		events
//	@Summary	Find event
//	@Produce	json
//	@Success	200	{array}		models.Event
//	@Failure	400	{object}	dtos.ResponseDTO
//	@Failure	500	{object}	dtos.ResponseDTO
//	@Router		/v1/events [get]
func (eventController *EventController) GetAll(c *gin.Context) {
	eventController.controller.GetAll(c)
}

//	@Tags		events
//	@Summary	Options
//	@Router		/v1/events [options]
func (eventController *EventController) Options(c *gin.Context) {
	eventController.controller.Options(c)
}
