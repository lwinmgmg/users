package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lwinmgmg/user/datamodels"
	"github.com/lwinmgmg/user/models"
	"gorm.io/gorm"
)

func (ctrl *UserController) GetUserByUserCode(ctx *gin.Context) {
	userCode := ctx.Param("userCode")
	var user models.User
	if _, err := user.GetPartnerByCode(userCode, ctrl.DB); err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.AbortWithStatusJSON(http.StatusNotFound, datamodels.DefaultResponse{
				Code:    1,
				Message: "User not found",
			})
			return
		}
		ctx.AbortWithStatusJSON(http.StatusBadRequest, datamodels.DefaultResponse{
			Code:    1,
			Message: "Failed to get partner!" + err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, datamodels.UserData{
		Code:            user.Code,
		IsAuthenticator: user.IsAuthenticator,
		Is2FA:           user.OtpUrl != "",
		PartnerData: datamodels.PartnerData{
			FirstName:        user.Partner.FirstName,
			LastName:         user.Partner.LastName,
			Email:            user.Partner.Email,
			Phone:            user.Partner.Phone,
			IsPhoneConfirmed: user.Partner.IsPhoneConfirmed,
			IsEmailConfirmed: user.Partner.IsEmailConfirmed,
			Code:             user.Partner.Code,
		},
	})
}
