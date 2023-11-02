package user

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lwinmgmg/user/datamodels"
	"github.com/lwinmgmg/user/middlewares"
	"github.com/lwinmgmg/user/models"
	"github.com/lwinmgmg/user/services"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func (ctrl *UserAuthController) OtpAuthenticate(ctx *gin.Context) {
	// Parse Input
	otpData := datamodels.OtpData{}
	if err := ctx.ShouldBindJSON(&otpData); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, datamodels.DefaultResponse{
			Code:    0,
			Message: fmt.Sprintf("Request must be json format [%v]", err.Error()),
		})
		return
	}
	// Get otp URL from redis
	val, err := services.GetKey(otpData.AccessToken)
	if err != nil {
		if err == redis.Nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, datamodels.DefaultResponse{
				Code:    1,
				Message: "Token is already expired",
			})
			return
		}
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, datamodels.DefaultResponse{
			Code:    1,
			Message: fmt.Sprintf("Authorization Required! [%v]", err),
		})
		return
	}
	// Parse Key from url
	valList := strings.Split(val, OtpKeyDivider)
	otpUrl := valList[0]
	userCode := valList[1]
	confirmType := valList[2]
	key, err := otp.NewKeyFromURL(otpUrl)
	fmt.Println("OtpAuth", key, key.Secret())
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, datamodels.DefaultResponse{
			Code:    1,
			Message: fmt.Sprintf("Internal Server Error. %v", err),
		})
		return
	}
	// Validate Passcode
	if totp.Validate(otpData.PassCode, key.Secret()) {
		var user models.User

		if err := ctrl.DB.Transaction(func(tx *gorm.DB) error {
			partner, err := user.GetPartnerByCode(userCode, tx)
			if err != nil {
				return err
			}
			switch confirmType {
			case string(OtpEmail):
				return partner.SetEmailConfirm(true, tx)
			case string(OtpPhone):
				return partner.SetPhoneConfirm(true, tx)
			case string(OtpAuthr):
				return user.SetIsAuthenticator(true, tx)
			case string(OtpEnable):
				return user.SetOtpUrl(otpUrl, tx)
			}
			return nil
		}); err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, datamodels.DefaultResponse{
				Code:    1,
				Message: fmt.Sprintf("Internal Server Error. %v", err),
			})
			return
		}
		services.DelKey([]string{otpData.AccessToken})
		ctx.JSON(http.StatusOK, datamodels.TokenResponse{
			AccessToken: middlewares.GetDefaultToken(userCode, string(user.Password), middlewares.DefaultTokenKey),
			TokenType:   middlewares.BearerTokenType,
		})
		return
	}
	fmt.Println(otpData.PassCode, key)
	ctx.AbortWithStatusJSON(http.StatusUnauthorized, datamodels.DefaultResponse{
		Code:    2,
		Message: "Authorization Required! [Invalid PassCode]",
	})
}
