package v1

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lwinmgmg/user/datamodels"
	"github.com/lwinmgmg/user/models"
	"github.com/lwinmgmg/user/services"
)

type UserController struct {
	Router *gin.RouterGroup
}

func (ctrl *UserController) HandleRoutes() {
	ctrl.Router.GET("/func/users/change_password", ctrl.ChangePassword)
	ctrl.Router.GET("/func/users/change_email", ctrl.ChangeEmail)
	ctrl.Router.GET("/func/users/change_phone", ctrl.ChangePhone)
}

func (ctrl *UserController) GenerateOtp(ctx *gin.Context) {
	username, ok := ctx.Get("username")
	userStr, ok1 := username.(string)
	if !ok || !ok1 {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, datamodels.DefaultResponse{
			Code:    1,
			Message: "Authorization Required!",
		})
		return
	}
	user := models.User{}
	_, err := user.GetPartnerByUsername(userStr, DB)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, datamodels.DefaultResponse{
			Code:    2,
			Message: "User not found",
		})
		return
	}
	if _, err := services.SetKey(fmt.Sprintf("otp_%v", userStr), "1", time.Minute); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, datamodels.DefaultResponse{
			Code:    1,
			Message: fmt.Sprintf("Can't set key %v", err),
		})
		return
	}
	ctx.JSON(http.StatusOK, datamodels.DefaultResponse{
		Code: 0,
		Message: "Success",
	})
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

func (ctrl *UserController) ConfirmEmail(ctx *gin.Context) {
	username, ok := ctx.Get("username")
	userStr, ok1 := username.(string)
	if !ok || !ok1 {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, datamodels.DefaultResponse{
			Code:    1,
			Message: "Authorization Required!",
		})
		return
	}
	user := models.User{}
	_, err := user.GetPartnerByUsername(userStr, DB)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, datamodels.DefaultResponse{
			Code:    2,
			Message: "User not found",
		})
		return
	}

}

func (ctrl *UserController) ChangePhone(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, map[string]string{
		"foo": "bar",
	})
}

func (ctrl *UserController) ConfirmPhone(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, map[string]string{
		"foo": "bar",
	})
}
