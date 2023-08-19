package v1

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/lwinmgmg/user/datamodels"
	"github.com/lwinmgmg/user/models"
	"github.com/lwinmgmg/user/services"
	"github.com/lwinmgmg/user/utils"
	"github.com/pquerna/otp/totp"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type UserAuthController struct {
	Router *gin.RouterGroup
}

func (ctrl *UserAuthController) HandleRoutes() {
	ctrl.Router.POST("/func/users/login", ctrl.Login)
	ctrl.Router.POST("/func/users/signup", ctrl.SignUp)
	ctrl.Router.POST("/func/users/otp_login", ctrl.OtpAuthenticate)
	ctrl.Router.POST("/func/users/re_auth", ctrl.ReAuthenticate)
}

func (ctrl *UserAuthController) Login(ctx *gin.Context) {
	userLoginData := datamodels.UserLoginData{}
	if err := ctx.ShouldBindJSON(&userLoginData); err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, datamodels.DefaultResponse{
			Code:    1,
			Message: fmt.Sprintf("Authorization Required! [%v]", err.Error()),
		})
		return
	}
	if err := userLoginData.Validate(); err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, datamodels.DefaultResponse{
			Code:    2,
			Message: fmt.Sprintf("Authorization Required! [%v]", err.Error()),
		})
		return
	}
	user := models.User{}
	if err := user.Authenticate(DB, &userLoginData); err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.AbortWithStatusJSON(http.StatusNotFound, datamodels.DefaultResponse{
				Code:    1,
				Message: fmt.Sprintf("User not found!"),
			})
			return
		}
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, datamodels.DefaultResponse{
			Code:    3,
			Message: fmt.Sprintf("Authorization Required! [%v]", err.Error()),
		})
		return
	}
	if user.Secret == "" {
		ctx.JSON(http.StatusOK, datamodels.TokenResponse{
			AccessToken: utils.GetDefaultToken(user.Username, utils.DefaultTokenKey),
			TokenType:   utils.BearerTokenType,
		})
		return
	}
	randomUuid := uuid.New()
	uuidString := randomUuid.String()
	partner := models.Partner{}
	if err := partner.GetPartnerByID(user.PartnerID, DB); err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.AbortWithStatusJSON(http.StatusNotFound, datamodels.DefaultResponse{
				Code:    2,
				Message: fmt.Sprintf("Partner not found!"),
			})
			return
		}
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, datamodels.DefaultResponse{
			Code:    1,
			Message: fmt.Sprintf("Internal Server Error. %v", err),
		})
		return
	}
	if _, err := services.SetKey(uuidString, fmt.Sprintf("%v:%v", user.Secret, user.Username), time.Minute); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, datamodels.DefaultResponse{
			Code:    1,
			Message: fmt.Sprintf("Internal Server ERROR : %v", err),
		})
		return
	}
	passCode, _ := totp.GenerateCode(user.Secret, time.Now().UTC().Add(-time.Minute))
	go services.MailSender.Send(passCode, []string{partner.Email})
	ctx.JSON(http.StatusAccepted, datamodels.TokenResponse{
		AccessToken: uuidString,
		TokenType:   utils.OtpTokenType,
	})
}

func (ctrl *UserAuthController) SignUp(ctx *gin.Context) {
	userSignUpData := datamodels.UserSignUpData{}
	if err := ctx.ShouldBindJSON(&userSignUpData); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, datamodels.DefaultResponse{
			Code:    0,
			Message: fmt.Sprintf("Request must be json format [%v]", err.Error()),
		})
		return
	}
	partner := models.Partner{
		FirstName: userSignUpData.FirstName,
		LastName:  userSignUpData.LastName,
		Email:     userSignUpData.Email,
		Phone:     userSignUpData.Phone,
	}
	user := models.User{
		Username: userSignUpData.UserName,
		Password: utils.Hash256(userSignUpData.Password),
	}
	if err := user.Exist(DB); err != nil {
		if err == utils.ErrRecordAlreadyExist {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, datamodels.DefaultResponse{
				Code:    1,
				Message: "User already exist",
			})
			return
		}
		ctx.AbortWithStatusJSON(http.StatusBadRequest, datamodels.DefaultResponse{
			Code:    1,
			Message: fmt.Sprintf("Error on user check [%v]", err.Error()),
		})
		return
	}
	if err := partner.CheckEmail(DB); err != nil {
		if err == utils.ErrRecordAlreadyExist {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, datamodels.DefaultResponse{
				Code:    2,
				Message: "Email already exist",
			})
			return
		}
		ctx.AbortWithStatusJSON(http.StatusBadRequest, datamodels.DefaultResponse{
			Code:    2,
			Message: fmt.Sprintf("Error on email check [%v]", err.Error()),
		})
		return
	}
	if err := partner.CheckPhone(DB); err != nil {
		if err == utils.ErrRecordAlreadyExist {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, datamodels.DefaultResponse{
				Code:    3,
				Message: "Phone already exist",
			})
			return
		}
		ctx.AbortWithStatusJSON(http.StatusBadRequest, datamodels.DefaultResponse{
			Code:    3,
			Message: fmt.Sprintf("Error on phone check [%v]", err.Error()),
		})
		return
	}
	if err := DB.Transaction(func(tx *gorm.DB) error {
		if err := partner.Create(tx); err != nil {
			return err
		}
		user.PartnerID = partner.ID
		if err := user.Create(tx); err != nil {
			return err
		}
		return nil
	}); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, datamodels.DefaultResponse{
			Code:    4,
			Message: fmt.Sprintf("Error occur on creating user [%v]", err.Error()),
		})
		return
	}
	user.Partner = partner
	ctx.JSON(http.StatusOK, datamodels.TokenResponse{
		AccessToken: utils.GetDefaultToken(user.Username, utils.DefaultTokenKey),
		TokenType:   utils.BearerTokenType,
	})
}

func (ctrl *UserAuthController) ReAuthenticate(ctx *gin.Context) {
	tokenData := datamodels.ReAuthTokenRequest{}
	if err := ctx.ShouldBindJSON(&tokenData); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, datamodels.DefaultResponse{
			Code:    0,
			Message: fmt.Sprintf("Request must be json format [%v]", err.Error()),
		})
		return
	}
	claim := jwt.RegisteredClaims{}
	if err := utils.ValidateToken(tokenData.Token, utils.DefaultTokenKey, &claim); err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			ctx.JSON(http.StatusOK, datamodels.TokenResponse{
				AccessToken: utils.GetDefaultToken(claim.ID, utils.DefaultTokenKey),
				TokenType:   utils.BearerTokenType,
			})
			return
		}
		ctx.AbortWithStatusJSON(http.StatusBadRequest, datamodels.DefaultResponse{
			Code:    0,
			Message: fmt.Sprintf("Token is invalid [%v]", err.Error()),
		})
		return
	}
	ctx.JSON(http.StatusAccepted, datamodels.TokenResponse{
		AccessToken: tokenData.Token,
		TokenType:   utils.BearerTokenType,
	})
}

func (ctrl *UserAuthController) OtpAuthenticate(ctx *gin.Context) {
	otpData := datamodels.OtpData{}
	if err := ctx.ShouldBindJSON(&otpData); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, datamodels.DefaultResponse{
			Code:    0,
			Message: fmt.Sprintf("Request must be json format [%v]", err.Error()),
		})
		return
	}
	val, err := services.GetKey(otpData.Token)
	valList := strings.Split(val, ":")
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
	if totp.Validate(otpData.PassCode, valList[0]) {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, datamodels.DefaultResponse{
			Code:    2,
			Message: "Authorization Required! [Invalid PassCode]",
		})
		return
	}
	ctx.JSON(http.StatusOK, datamodels.TokenResponse{
		AccessToken: utils.GetDefaultToken(valList[1], utils.DefaultTokenKey),
		TokenType:   utils.BearerTokenType,
	})
}
