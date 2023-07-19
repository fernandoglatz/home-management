package api

import (
	"github.com/fernandoglatz/home-management/backend/api/dtos"
	"github.com/fernandoglatz/home-management/backend/models"
	"github.com/gin-gonic/gin"
)

type DeviceController struct {
	controller Controller[*models.Device, *dtos.DeviceDTO]
}

func NewDeviceController() DeviceController {
	return DeviceController{
		controller: NewController[*models.Device, *dtos.DeviceDTO]("devices"),
	}
}

//	@Tags		devices
//	@Summary	Create device
//	@Param		request	body	dtos.DeviceDTO	true	"body"
//	@Accept		json
//	@Produce	json
//	@Success	200	{object}	models.Device
//	@Failure	400	{object}	dtos.ResponseDTO
//	@Failure	500	{object}	dtos.ResponseDTO
//	@Router		/v1/devices [post]
func (deviceController *DeviceController) Create(c *gin.Context) {
	deviceController.controller.Create(c)
}

//	@Tags		devices
//	@Summary	Update device
//	@Param		id		path	string			true	"Device ID"
//	@Param		request	body	dtos.DeviceDTO	true	"body"
//	@Accept		json
//	@Produce	json
//	@Success	200	{object}	models.Device
//	@Failure	400	{object}	dtos.ResponseDTO
//	@Failure	500	{object}	dtos.ResponseDTO
//	@Router		/v1/devices/{id} [put]
func (deviceController *DeviceController) Update(c *gin.Context) {
	deviceController.controller.Update(c)
}

//	@Tags		devices
//	@Summary	Partial update device
//	@Param		id		path	string			true	"Device ID"
//	@Param		request	body	dtos.DeviceDTO	true	"body"
//	@Accept		json
//	@Produce	json
//	@Success	200	{object}	models.Device
//	@Failure	400	{object}	dtos.ResponseDTO
//	@Failure	500	{object}	dtos.ResponseDTO
//	@Router		/v1/devices/{id} [patch]
func (deviceController *DeviceController) Patch(c *gin.Context) {
	deviceController.controller.Patch(c)
}

//	@Tags		devices
//	@Summary	Delete device
//	@Param		id	path	string	true	"Device ID"
//	@Produce	json
//	@Failure	200	{object}	dtos.ResponseDTO
//	@Failure	400	{object}	dtos.ResponseDTO
//	@Failure	500	{object}	dtos.ResponseDTO
//	@Router		/v1/devices/{id} [delete]
func (deviceController *DeviceController) Delete(c *gin.Context) {
	deviceController.controller.Delete(c)
}

//	@Tags		devices
//	@Summary	Check if device exists
//	@Param		id	path	string	true	"Device ID"
//	@Router		/v1/devices/{id} [head]
func (deviceController *DeviceController) Head(c *gin.Context) {
	deviceController.controller.Head(c)
}

//	@Tags		devices
//	@Summary	Find device
//	@Param		id	path	string	true	"Device ID"
//	@Produce	json
//	@Success	200	{object}	models.Device
//	@Failure	400	{object}	dtos.ResponseDTO
//	@Failure	500	{object}	dtos.ResponseDTO
//	@Router		/v1/devices/{id} [get]
func (deviceController *DeviceController) Get(c *gin.Context) {
	deviceController.controller.Get(c)
}

//	@Tags		devices
//	@Summary	Find device
//	@Produce	json
//	@Success	200	{array}		models.Device
//	@Failure	400	{object}	dtos.ResponseDTO
//	@Failure	500	{object}	dtos.ResponseDTO
//	@Router		/v1/devices [get]
func (deviceController *DeviceController) GetAll(c *gin.Context) {
	deviceController.controller.GetAll(c)
}

//	@Tags		devices
//	@Summary	Options
//	@Router		/v1/devices [options]
func (deviceController *DeviceController) Options(c *gin.Context) {
	deviceController.controller.Options(c)
}
