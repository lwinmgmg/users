package v1

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lwinmgmg/user/datamodels"
	"github.com/lwinmgmg/user/models"
	"github.com/lwinmgmg/user/services"
	"github.com/lwinmgmg/user/utils"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"gorm.io/gorm"
)

type UserController struct {
	Router *gin.RouterGroup
}

func (ctrl *UserController) HandleRoutes() {
	ctrl.Router.GET("/func/users/enable_two_factor", ctrl.EnableTwoFactorAuth)
	ctrl.Router.GET("/func/users/confirm_email", ctrl.ConfirmEmail)
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
		Code:    0,
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
	username, ok := GetUserFromContext(ctx)
	if !ok {
		return
	}
	user := models.User{}
	// Get Partner
	partner, err := user.GetPartnerByUsername(username, DB)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, datamodels.DefaultResponse{
			Code:    2,
			Message: "Partner not found",
		})
		return
	}
	// Check email is already confirmed or not
	if partner.IsEmailConfirmed {
		ctx.JSON(http.StatusAccepted, datamodels.DefaultResponse{
			Code:    1,
			Message: "Email is already confirmed",
		})
	}
	// Generate UUID
	randomUuid := uuid.New()
	uuidString := randomUuid.String()
	tokenExpireTime := 5 * time.Minute
	otpUrl, err := utils.GenerateOtpUrl(username, tokenExpireTime)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, datamodels.DefaultResponse{
			Code:    2,
			Message: fmt.Sprintf("Internal Server ERROR : %v", err),
		})
		return
	}
	if _, err := services.SetKey(uuidString, fmt.Sprintf(OPT_UUID_FORMAT, otpUrl, user.Username, OtpEmail), tokenExpireTime); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, datamodels.DefaultResponse{
			Code:    2,
			Message: fmt.Sprintf("Internal Server ERROR : %v", err),
		})
		return
	}
	// Parse Key from url
	key, err := otp.NewKeyFromURL(otpUrl)
	fmt.Println("Confirm email", key, key.Secret())
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, datamodels.DefaultResponse{
			Code:    1,
			Message: fmt.Sprintf("Internal Server Error. %v", err),
		})
		return
	}
	// Get passcode and send Email
	passCode, _ := totp.GenerateCode(key.Secret(), time.Now().UTC())
	go services.MailSender.Send(passCode, []string{partner.Email})
	ctx.JSON(http.StatusOK, datamodels.TokenResponse{
		AccessToken: uuidString,
		TokenType:   utils.OtpTokenType,
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

func (ctrl *UserController) EnableTwoFactorAuth(ctx *gin.Context) {
	username, ok := GetUserFromContext(ctx)
	if !ok {
		return
	}
	var user models.User
	if err := DB.Transaction(func(tx *gorm.DB) error {
		// Get user
		if err := user.GetUserByUsername(username, tx); err != nil {
			return err
		}
		partner, err := user.GetPartnerByUsername(user.Username, tx)
		if err != nil {
			return err
		}
		if !(partner.IsEmailConfirmed || partner.IsPhoneConfirmed) {
			return utils.ErrInvalid
		}
		// Generate OTP URL
		otpUrl, err := utils.GenerateOtpUrl(user.Username, time.Minute)
		if err != nil {
			return err
		}
		return user.SetOtpUrl(otpUrl, tx)
	}); err != nil {
		if err == utils.ErrInvalid {
			ctx.JSON(http.StatusAccepted, datamodels.DefaultResponse{
				Code:    1,
				Message: "Confirm Email Or Phone first",
			})
			return
		}
		ctx.AbortWithStatusJSON(http.StatusBadRequest, datamodels.DefaultResponse{
			Code:    1,
			Message: fmt.Sprintf("Two Factor Authentication can't be set. [%v]", err),
		})
		return
	}
	ctx.JSON(http.StatusOK, datamodels.DefaultResponse{
		Code:    1,
		Message: "Successfully enabled two factor authentication",
	})
}
