package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lwinmgmg/user/datamodels"
	"github.com/lwinmgmg/user/services"
	"gorm.io/gorm"
)

var (
	DB *gorm.DB = services.PgDb
)

func GetUserFromContext(ctx *gin.Context) (string, bool) {
	username, ok := ctx.Get("username")
	userStr, ok1 := username.(string)
	if !ok || !ok1 {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, datamodels.DefaultResponse{
			Code:    1,
			Message: "Authorization Required!",
		})
		return "", false
	}
	return userStr, true
}
