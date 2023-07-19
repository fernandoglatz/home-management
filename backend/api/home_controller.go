package api

import (
	"github.com/fernandoglatz/home-management/backend/api/dtos"
	"github.com/fernandoglatz/home-management/backend/models"
	"github.com/gin-gonic/gin"
)

type HomeController struct {
	controller Controller[*models.Home, *dtos.HomeDTO]
}

func NewHomeController() HomeController {
	return HomeController{
		controller: NewController[*models.Home, *dtos.HomeDTO]("homes"),
	}
}

//	@Tags		homes
//	@Summary	Create home
//	@Param		request	body	dtos.HomeDTO	true	"body"
//	@Accept		json
//	@Produce	json
//	@Success	200	{object}	models.Home
//	@Failure	400	{object}	dtos.ResponseDTO
//	@Failure	500	{object}	dtos.ResponseDTO
//	@Router		/v1/homes [post]
func (homeController *HomeController) Create(c *gin.Context) {
	homeController.controller.Create(c)
}

//	@Tags		homes
//	@Summary	Update home
//	@Param		id		path	string			true	"Home ID"
//	@Param		request	body	dtos.HomeDTO	true	"body"
//	@Accept		json
//	@Produce	json
//	@Success	200	{object}	models.Home
//	@Failure	400	{object}	dtos.ResponseDTO
//	@Failure	500	{object}	dtos.ResponseDTO
//	@Router		/v1/homes/{id} [put]
func (homeController *HomeController) Update(c *gin.Context) {
	homeController.controller.Update(c)
}

//	@Tags		homes
//	@Summary	Partial update home
//	@Param		id		path	string			true	"Home ID"
//	@Param		request	body	dtos.HomeDTO	true	"body"
//	@Accept		json
//	@Produce	json
//	@Success	200	{object}	models.Home
//	@Failure	400	{object}	dtos.ResponseDTO
//	@Failure	500	{object}	dtos.ResponseDTO
//	@Router		/v1/homes/{id} [patch]
func (homeController *HomeController) Patch(c *gin.Context) {
	homeController.controller.Patch(c)
}

//	@Tags		homes
//	@Summary	Delete home
//	@Param		id	path	string	true	"Home ID"
//	@Produce	json
//	@Failure	200	{object}	dtos.ResponseDTO
//	@Failure	400	{object}	dtos.ResponseDTO
//	@Failure	500	{object}	dtos.ResponseDTO
//	@Router		/v1/homes/{id} [delete]
func (homeController *HomeController) Delete(c *gin.Context) {
	homeController.controller.Delete(c)
}

//	@Tags		homes
//	@Summary	Check if home exists
//	@Param		id	path	string	true	"Home ID"
//	@Router		/v1/homes/{id} [head]
func (homeController *HomeController) Head(c *gin.Context) {
	homeController.controller.Head(c)
}

//	@Tags		homes
//	@Summary	Find home
//	@Param		id	path	string	true	"Home ID"
//	@Produce	json
//	@Success	200	{object}	models.Home
//	@Failure	400	{object}	dtos.ResponseDTO
//	@Failure	500	{object}	dtos.ResponseDTO
//	@Router		/v1/homes/{id} [get]
func (homeController *HomeController) Get(c *gin.Context) {
	homeController.controller.Get(c)
}

//	@Tags		homes
//	@Summary	Find home
//	@Produce	json
//	@Success	200	{array}		models.Home
//	@Failure	400	{object}	dtos.ResponseDTO
//	@Failure	500	{object}	dtos.ResponseDTO
//	@Router		/v1/homes [get]
func (homeController *HomeController) GetAll(c *gin.Context) {
	homeController.controller.GetAll(c)
}

//	@Tags		homes
//	@Summary	Options
//	@Router		/v1/homes [options]
func (homeController *HomeController) Options(c *gin.Context) {
	homeController.controller.Options(c)
}
