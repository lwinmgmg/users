package user

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image/png"
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
)

func (ctrl *UserController) EnableAuthenticator(ctx *gin.Context) {
	userCode, ok := controllers.GetUserFromContext(ctx)
	if !ok {
		return
	}
	var user models.User

	if err := user.GetUserByCode(userCode, ctrl.DB); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, datamodels.DefaultResponse{
			Code:    1,
			Message: fmt.Sprintf("Can't set authenticator [%v]", err),
		})
		return
	}
	if user.OtpUrl == "" {
		ctx.JSON(http.StatusBadRequest, datamodels.DefaultResponse{
			Code:    1,
			Message: "Enable two factor authentication first",
		})
		return
	}
	randomUuid := uuid.New()
	uuidString := randomUuid.String()
	tokenExpireTime := 5 * time.Minute
	if _, err := services.SetKey(uuidString, fmt.Sprintf(OPT_UUID_FORMAT, user.OtpUrl, user.Code, OtpAuthr), tokenExpireTime); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, datamodels.DefaultResponse{
			Code:    2,
			Message: fmt.Sprintf("Internal Server ERROR : %v", err),
		})
		return
	}
	var buf bytes.Buffer
	key, err := otp.NewKeyFromURL(user.OtpUrl)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, datamodels.DefaultResponse{
			Code:    1,
			Message: fmt.Sprintf("Internal Server Error. %v", err),
		})
		return
	}
	img, err := key.Image(100, 100)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, datamodels.DefaultResponse{
			Code:    4,
			Message: fmt.Sprintf("Internal Server Error. %v", err),
		})
		return
	}
	if err := png.Encode(&buf, img); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, datamodels.DefaultResponse{
			Code:    1,
			Message: fmt.Sprintf("Internal Server Error. %v", err),
		})
		return
	}
	ctx.JSON(http.StatusOK, datamodels.TokenAuthResponse{
		AccessToken: uuidString,
		TokenType:   middlewares.OtpTokenType,
		Image:       base64.StdEncoding.EncodeToString(buf.Bytes()),
		Key:         key.Secret(),
	})
}
