package user

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/lwinmgmg/user/datamodels"
	"github.com/lwinmgmg/user/middlewares"
	"github.com/lwinmgmg/user/models"
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
	if err := middlewares.ValidateToken(tokenData.AccessToken, middlewares.DefaultTokenKey, &claim); err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			password, err := models.GetPasswordByUserCode(claim.Subject, ctrl.DB)
			if err != nil {
				return
			}
			ctx.JSON(http.StatusOK, datamodels.TokenResponse{
				AccessToken: middlewares.GetDefaultToken(claim.Subject, password, middlewares.DefaultTokenKey),
				TokenType:   middlewares.BearerTokenType,
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
		TokenType:   middlewares.BearerTokenType,
	})
}
