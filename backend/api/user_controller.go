package api

import (
	"github.com/fernandoglatz/home-management/backend/api/dtos"
	"github.com/fernandoglatz/home-management/backend/models"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	controller Controller[*models.User, *dtos.UserDTO]
}

func NewUserController() UserController {
	return UserController{
		controller: NewController[*models.User, *dtos.UserDTO]("users"),
	}
}

//	@Tags		users
//	@Summary	Create user
//	@Param		request	body	dtos.UserDTO	true	"body"
//	@Accept		json
//	@Produce	json
//	@Success	200	{object}	models.User
//	@Failure	400	{object}	dtos.ResponseDTO
//	@Failure	500	{object}	dtos.ResponseDTO
//	@Router		/v1/users [post]
func (userController *UserController) Create(c *gin.Context) {
	userController.controller.Create(c)
}

//	@Tags		users
//	@Summary	Update user
//	@Param		id		path	string			true	"User ID"
//	@Param		request	body	dtos.UserDTO	true	"body"
//	@Accept		json
//	@Produce	json
//	@Success	200	{object}	models.User
//	@Failure	400	{object}	dtos.ResponseDTO
//	@Failure	500	{object}	dtos.ResponseDTO
//	@Router		/v1/users/{id} [put]
func (userController *UserController) Update(c *gin.Context) {
	userController.controller.Update(c)
}

//	@Tags		users
//	@Summary	Partial update user
//	@Param		id		path	string			true	"User ID"
//	@Param		request	body	dtos.UserDTO	true	"body"
//	@Accept		json
//	@Produce	json
//	@Success	200	{object}	models.User
//	@Failure	400	{object}	dtos.ResponseDTO
//	@Failure	500	{object}	dtos.ResponseDTO
//	@Router		/v1/users/{id} [patch]
func (userController *UserController) Patch(c *gin.Context) {
	userController.controller.Patch(c)
}

//	@Tags		users
//	@Summary	Delete user
//	@Param		id	path	string	true	"User ID"
//	@Produce	json
//	@Failure	200	{object}	dtos.ResponseDTO
//	@Failure	400	{object}	dtos.ResponseDTO
//	@Failure	500	{object}	dtos.ResponseDTO
//	@Router		/v1/users/{id} [delete]
func (userController *UserController) Delete(c *gin.Context) {
	userController.controller.Delete(c)
}

//	@Tags		users
//	@Summary	Check if user exists
//	@Param		id	path	string	true	"User ID"
//	@Router		/v1/users/{id} [head]
func (userController *UserController) Head(c *gin.Context) {
	userController.controller.Head(c)
}

//	@Tags		users
//	@Summary	Find user
//	@Param		id	path	string	true	"User ID"
//	@Produce	json
//	@Success	200	{object}	models.User
//	@Failure	400	{object}	dtos.ResponseDTO
//	@Failure	500	{object}	dtos.ResponseDTO
//	@Router		/v1/users/{id} [get]
func (userController *UserController) Get(c *gin.Context) {
	userController.controller.Get(c)
}

//	@Tags		users
//	@Summary	Find user
//	@Produce	json
//	@Success	200	{array}		models.User
//	@Failure	400	{object}	dtos.ResponseDTO
//	@Failure	500	{object}	dtos.ResponseDTO
//	@Router		/v1/users [get]
func (userController *UserController) GetAll(c *gin.Context) {
	userController.controller.GetAll(c)
}

//	@Tags		users
//	@Summary	Options
//	@Router		/v1/users [options]
func (userController *UserController) Options(c *gin.Context) {
	userController.controller.Options(c)
}
