package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lwinmgmg/user/datamodels"
)

func GetUserFromContext(ctx *gin.Context) (string, bool) {
	userCode, ok := ctx.Get("userCode")
	userCodeStr, ok1 := userCode.(string)
	if !ok || !ok1 {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, datamodels.DefaultResponse{
			Code:    1,
			Message: "Authorization Required!",
		})
		return "", false
	}
	return userCodeStr, true
}

type HttpController interface {
	HandleRoutes()
}
