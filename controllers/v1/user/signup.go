package user

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lwinmgmg/user/datamodels"
	"github.com/lwinmgmg/user/middlewares"
	"github.com/lwinmgmg/user/models"
	"github.com/lwinmgmg/user/utils"
	"gorm.io/gorm"
)

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
	err := user.GetUserByUsername(user.Username, ctrl.DB)
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
	if err := partner.CheckEmail(ctrl.DB); err != nil {
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
	if err := partner.CheckPhone(ctrl.DB); err != nil {
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
	if err := ctrl.DB.Transaction(func(tx *gorm.DB) error {
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
		AccessToken: middlewares.GetDefaultToken(user.Code, string(user.Password), middlewares.DefaultTokenKey),
		TokenType:   middlewares.BearerTokenType,
	})
}
