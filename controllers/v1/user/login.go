package user

import (
	"errors"
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
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"gorm.io/gorm"
)

type UserLoginData struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (userLoginData *UserLoginData) Validate() error {
	if userLoginData.Username == "" {
		return errors.New("wrong username")
	}
	if userLoginData.Password == "" {
		return errors.New("wrong password")
	}
	return nil
}

func (ctrl *UserAuthController) Login(ctx *gin.Context) {
	userLoginData := UserLoginData{}
	if err := ctx.ShouldBindJSON(&userLoginData); err != nil {
		panic(controllers.NewPanicResponse(http.StatusUnauthorized, 1, fmt.Sprintf("Authorization Required! [%v]", err.Error())))
	}
	if err := userLoginData.Validate(); err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, datamodels.DefaultResponse{
			Code:    2,
			Message: fmt.Sprintf("Authorization Required! [%v]", err.Error()),
		})
		return
	}
	user := models.User{}
	if err := user.Authenticate(ctrl.DB, userLoginData.Username, userLoginData.Password); err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.AbortWithStatusJSON(http.StatusNotFound, datamodels.DefaultResponse{
				Code:    1,
				Message: "User not found!",
			})
			return
		}
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, datamodels.DefaultResponse{
			Code:    3,
			Message: fmt.Sprintf("Authorization Required! [%v]", err.Error()),
		})
		return
	}
	if user.OtpUrl == "" {
		ctx.JSON(http.StatusOK, datamodels.TokenResponse{
			AccessToken: middlewares.GetDefaultToken(user.Code, string(user.Password), middlewares.DefaultTokenKey),
			TokenType:   middlewares.BearerTokenType,
		})
		return
	}
	randomUuid := uuid.New()
	uuidString := randomUuid.String()
	tokenExpireTime := 5 * time.Minute
	if user.IsAuthenticator {
		tokenExpireTime = 5 * time.Minute
	}
	if _, err := services.SetKey(uuidString, fmt.Sprintf(OPT_UUID_FORMAT, user.OtpUrl, user.Code, OtpLogin), tokenExpireTime); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, datamodels.DefaultResponse{
			Code:    2,
			Message: fmt.Sprintf("Internal Server ERROR : %v", err),
		})
		return
	}
	if user.IsAuthenticator {
		ctx.JSON(http.StatusCreated, datamodels.TokenResponse{
			AccessToken: uuidString,
			TokenType:   middlewares.OtpTokenType,
		})
		return
	}
	partner := models.Partner{}
	if err := partner.GetPartnerByID(user.PartnerID, ctrl.DB); err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.AbortWithStatusJSON(http.StatusNotFound, datamodels.DefaultResponse{
				Code:    2,
				Message: "Partner not found!",
			})
			return
		}
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, datamodels.DefaultResponse{
			Code:    1,
			Message: fmt.Sprintf("Internal Server Error. %v", err),
		})
		return
	}
	key, err := otp.NewKeyFromURL(user.OtpUrl)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, datamodels.DefaultResponse{
			Code:    1,
			Message: fmt.Sprintf("Internal Server Error. %v", err),
		})
		return
	}
	passCode, _ := totp.GenerateCode(key.Secret(), time.Now().UTC())
	go services.MailSender.Send(passCode, []string{partner.Email})
	ctx.JSON(http.StatusCreated, datamodels.TokenResponse{
		AccessToken: uuidString,
		TokenType:   middlewares.OtpTokenType,
	})
}
