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

func (ctrl *UserController) ConfirmEmail(ctx *gin.Context) {
	userCode, ok := controllers.GetUserFromContext(ctx)
	if !ok {
		return
	}
	user := models.User{}
	// Get Partner
	partner, err := user.GetPartnerByCode(userCode, ctrl.DB)
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
		return
	}
	// Generate UUID
	randomUuid := uuid.New()
	uuidString := randomUuid.String()
	tokenExpireTime := 5 * time.Minute
	otpUrl, err := utils.GenerateOtpUrl(user.Username, tokenExpireTime)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, datamodels.DefaultResponse{
			Code:    2,
			Message: fmt.Sprintf("Internal Server ERROR : %v", err),
		})
		return
	}
	if _, err := services.SetKey(uuidString, fmt.Sprintf(OPT_UUID_FORMAT, otpUrl, user.Code, OtpEmail), tokenExpireTime); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, datamodels.DefaultResponse{
			Code:    2,
			Message: fmt.Sprintf("Internal Server ERROR : %v", err),
		})
		return
	}
	// Parse Key from url
	key, err := otp.NewKeyFromURL(otpUrl)
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
		TokenType:   middlewares.OtpTokenType,
	})
}
