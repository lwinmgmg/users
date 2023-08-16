package v1

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/lwinmgmg/user/datamodels"
	"github.com/lwinmgmg/user/models"
	"github.com/lwinmgmg/user/utils"
	"gorm.io/gorm"
)

type UserAuthController struct {
	Router *gin.RouterGroup
}

func (ctrl *UserAuthController) HandleRoutes() {
	ctrl.Router.POST("/func/users/login", ctrl.Login)
	ctrl.Router.POST("/func/users/signup", ctrl.SignUp)
	ctrl.Router.POST("/func/users/re_auth", ctrl.ReAuthenticate)
}

func (ctrl *UserAuthController) Login(ctx *gin.Context) {
	userLoginData := datamodels.UserLoginData{}
	if err := ctx.ShouldBindJSON(&userLoginData); err != nil {
		ctx.JSON(http.StatusUnauthorized, datamodels.DefaultResponse{
			Code:    1,
			Message: fmt.Sprintf("Authorization Required! [%v]", err.Error()),
		})
		return
	}
	if err := userLoginData.Validate(); err != nil {
		ctx.JSON(http.StatusUnauthorized, datamodels.DefaultResponse{
			Code:    2,
			Message: fmt.Sprintf("Authorization Required! [%v]", err.Error()),
		})
		return
	}
	user := models.User{}
	if err := user.Authenticate(DB, &userLoginData); err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, datamodels.DefaultResponse{
				Code:    1,
				Message: fmt.Sprintf("User not found!"),
			})
			return
		}
		ctx.JSON(http.StatusUnauthorized, datamodels.DefaultResponse{
			Code:    3,
			Message: fmt.Sprintf("Authorization Required! [%v]", err.Error()),
		})
		return
	}
	ctx.JSON(http.StatusOK, datamodels.TokenResponse{
		AccessToken: utils.GetDefaultToken(user.Username, utils.DefaultTokenKey),
		TokenType:   utils.TokenType,
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
	ctx.JSON(http.StatusOK, user)
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
				TokenType:   utils.TokenType,
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
		TokenType:   utils.TokenType,
	})
}
