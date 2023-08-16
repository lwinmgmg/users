package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	Router *gin.RouterGroup
}

func (ctrl *UserController) HandleRoutes() {
	ctrl.Router.GET("/func/users/change_password", ctrl.ChangePassword)
	ctrl.Router.GET("/func/users/change_email", ctrl.ChangeEmail)
	ctrl.Router.GET("/func/users/change_phone", ctrl.ChangePhone)
}

func (ctrl *UserController) ChangePassword(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, map[string]string{
		"foo": "bar",
	})
}

func (ctrl *UserController) ChangeEmail(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, map[string]string{
		"foo": "bar",
	})
}

func (ctrl *UserController) ConfirmEmail(ctx *gin.Context){
	ctx.JSON(http.StatusOK, map[string]string{
		"foo": "bar",
	})
}

func (ctrl *UserController) ChangePhone(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, map[string]string{
		"foo": "bar",
	})
}

func (ctrl *UserController) ConfirmPhone(ctx *gin.Context){
	ctx.JSON(http.StatusOK, map[string]string{
		"foo": "bar",
	})
}
