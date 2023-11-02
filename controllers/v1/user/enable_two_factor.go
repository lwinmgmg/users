package user

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lwinmgmg/user/controllers"
	"github.com/lwinmgmg/user/datamodels"
	"github.com/lwinmgmg/user/middlewares"
	"github.com/lwinmgmg/user/models"
	"github.com/lwinmgmg/user/services"
	"github.com/lwinmgmg/user/utils"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

func (ctrl *UserController) EnableTwoFactorAuth(ctx *gin.Context) {
	userCode, ok := controllers.GetUserFromContext(ctx)
	if !ok {
		return
	}
	var user models.User
	partner, err := user.GetPartnerByCode(userCode, ctrl.DB)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, datamodels.DefaultResponse{
			Code:    1,
			Message: fmt.Sprintf("Two Factor Authentication can't be set. [%v]", err),
		})
	}
	if !(partner.IsEmailConfirmed || partner.IsPhoneConfirmed) {
		ctx.JSON(http.StatusAccepted, datamodels.DefaultResponse{
			Code:    1,
			Message: "Confirm Email Or Phone first",
		})
		return
	}
	// Generate OTP URL
	otpUrl, err := utils.GenerateOtpUrl(user.Username, time.Minute)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, datamodels.DefaultResponse{
			Code:    1,
			Message: fmt.Sprintf("Two Factor Authentication can't be set. [%v]", err),
		})
	}
	randomUuid := uuid.New()
	uuidString := randomUuid.String()
	tokenExpireTime := 5 * time.Minute
	if _, err := services.SetKey(uuidString, fmt.Sprintf(OPT_UUID_FORMAT, otpUrl, user.Code, OtpEnable), tokenExpireTime); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, datamodels.DefaultResponse{
			Code:    2,
			Message: fmt.Sprintf("Internal Server ERROR : %v", err),
		})
		return
	}

	key, err := otp.NewKeyFromURL(otpUrl)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, datamodels.DefaultResponse{
			Code:    1,
			Message: fmt.Sprintf("Internal Server Error. %v", err),
		})
		return
	}
	passCode, _ := totp.GenerateCode(key.Secret(), time.Now().UTC())
	go services.MailSender.Send(passCode, []string{partner.Email})
	ctx.JSON(http.StatusOK, datamodels.TokenResponse{
		AccessToken: uuidString,
		TokenType:   middlewares.OtpTokenType,
	})
}
