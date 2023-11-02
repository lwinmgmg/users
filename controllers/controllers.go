package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lwinmgmg/user/datamodels"
)

type DefaultResponse struct {
	Code    int            `json:"code,omitempty"`
	Message string         `json:"message,omitempty"`
	Data    map[string]any `json:"data,omitempty"`
}

type PanicResponse struct {
	Response   DefaultResponse
	HttpStatus int
}

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

func CusRecoveryFunction(ctx *gin.Context, err any) {
	switch err := err.(type) {
	case PanicResponse:
		ctx.AbortWithStatusJSON(err.HttpStatus, err.Response)
		return
	default:
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, DefaultResponse{
			Code:    500,
			Message: fmt.Sprintf("Unknown error, %v", err),
		})
	}
}

func NewPanicResponse(httpStatus, code int, mesg string, data ...map[string]any) PanicResponse {
	respData := map[string]any{}
	if len(data) > 0 {
		respData = data[0]
	}
	return PanicResponse{
		HttpStatus: httpStatus,
		Response: DefaultResponse{
			Code:    code,
			Message: mesg,
			Data:    respData,
		},
	}
}
