package user

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var (
	OPT_UUID_FORMAT string = fmt.Sprintf("%v%v%v%v%v", "%v", OtpKeyDivider, "%v", OtpKeyDivider, "%v") //otp url, username, type
)

type UserAuthController struct {
	Router gin.IRoutes
	DB     *gorm.DB
}

func (ctrl *UserAuthController) HandleRoutes() {
	ctrl.Router.POST("/func/users/login", ctrl.Login)
	ctrl.Router.POST("/func/users/signup", ctrl.SignUp)
	ctrl.Router.POST("/func/users/otp_login", ctrl.OtpAuthenticate)
	ctrl.Router.POST("/func/users/re_auth", ctrl.ReAuthenticate)
}

type UserController struct {
	Router *gin.RouterGroup
	DB     *gorm.DB
}

func (ctrl *UserController) HandleRoutes() {
	ctrl.Router.GET("/func/users/enable_two_factor", ctrl.EnableTwoFactorAuth)
	ctrl.Router.GET("/func/users/enable_auth", ctrl.EnableAuthenticator)
	ctrl.Router.GET("/func/users/confirm_email", ctrl.ConfirmEmail)
	ctrl.Router.GET("/func/users/change_password", ctrl.ChangePassword)
	ctrl.Router.GET("/func/users/change_email", ctrl.ChangeEmail)
	ctrl.Router.GET("/func/users/change_phone", ctrl.ChangePhone)
	ctrl.Router.GET("/users/code/:userCode", ctrl.GetUserByUserCode)
	ctrl.Router.GET("/users/profile", ctrl.GetProfile)
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
