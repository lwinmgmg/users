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
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type OtpConfirmType string

const (
	OtpKeyDivider                = "|||"
	OtpLogin      OtpConfirmType = "1"
	OtpEmail      OtpConfirmType = "2"
	OtpPhone      OtpConfirmType = "3"
	OtpAuthr      OtpConfirmType = "4"
	OtpEnable     OtpConfirmType = "5"
)

var (
	OPT_UUID_FORMAT string = fmt.Sprintf("%v%v%v%v%v", "%v", OtpKeyDivider, "%v", OtpKeyDivider, "%v") //otp url, username, type
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
	if user.OtpUrl == "" {
		ctx.JSON(http.StatusOK, datamodels.TokenResponse{
			AccessToken: utils.GetDefaultToken(user.Username, utils.DefaultTokenKey),
			TokenType:   utils.BearerTokenType,
		})
		return
	}
	randomUuid := uuid.New()
	uuidString := randomUuid.String()
	tokenExpireTime := 5 * time.Minute
	if user.IsAuthenticator {
		tokenExpireTime = 5 * time.Minute
	}
	if _, err := services.SetKey(uuidString, fmt.Sprintf(OPT_UUID_FORMAT, user.OtpUrl, user.Username, OtpLogin), tokenExpireTime); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, datamodels.DefaultResponse{
			Code:    2,
			Message: fmt.Sprintf("Internal Server ERROR : %v", err),
		})
		return
	}
	if user.IsAuthenticator {
		ctx.JSON(http.StatusAccepted, datamodels.TokenResponse{
			AccessToken: uuidString,
			TokenType:   utils.OtpTokenType,
		})
		return
	}
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
	user := models.User{
		Username: userSignUpData.UserName,
		Password: utils.Hash256(userSignUpData.Password),
	}
	err := user.GetUserByUsername(user.Username, DB)
	if err == nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, datamodels.DefaultResponse{
			Code:    1,
			Message: "User already exists",
		})
		return
	}
	if err != gorm.ErrRecordNotFound {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, datamodels.DefaultResponse{
			Code:    1,
			Message: fmt.Sprintf("Internal server error [%v]", err.Error()),
		})
		return
	}
	partner := models.Partner{
		FirstName: userSignUpData.FirstName,
		LastName:  userSignUpData.LastName,
		Email:     userSignUpData.Email,
		Phone:     userSignUpData.Phone,
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
			Code:    3,
			Message: fmt.Sprintf("Error on email check [%v]", err.Error()),
		})
		return
	}
	if err := partner.CheckPhone(DB); err != nil {
		if err == utils.ErrRecordAlreadyExist {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, datamodels.DefaultResponse{
				Code:    4,
				Message: "Phone already exist",
			})
			return
		}
		ctx.AbortWithStatusJSON(http.StatusBadRequest, datamodels.DefaultResponse{
			Code:    5,
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
			Code:    6,
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
	// Parse Input
	tokenData := datamodels.ReAuthTokenRequest{}
	if err := ctx.ShouldBindJSON(&tokenData); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, datamodels.DefaultResponse{
			Code:    0,
			Message: fmt.Sprintf("Request must be json format [%v]", err.Error()),
		})
		return
	}
	// Parse claim and Validate token
	claim := jwt.RegisteredClaims{}
	if err := utils.ValidateToken(tokenData.AccessToken, utils.DefaultTokenKey, &claim); err != nil {
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
	// Token is not expire yet and send back the old one
	ctx.JSON(http.StatusAccepted, datamodels.TokenResponse{
		AccessToken: tokenData.AccessToken,
		TokenType:   utils.BearerTokenType,
	})
}

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
	username := valList[1]
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

		if err := DB.Transaction(func(tx *gorm.DB) error {
			partner, err := user.GetPartnerByUsername(username, tx)
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
			AccessToken: utils.GetDefaultToken(username, utils.DefaultTokenKey),
			TokenType:   utils.BearerTokenType,
		})
		return
	}
	fmt.Println(otpData.PassCode, key)
	ctx.AbortWithStatusJSON(http.StatusUnauthorized, datamodels.DefaultResponse{
		Code:    2,
		Message: "Authorization Required! [Invalid PassCode]",
	})

}
