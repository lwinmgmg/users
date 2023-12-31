package middlewares

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/lwinmgmg/user/datamodels"
	"github.com/lwinmgmg/user/utils"
)

func ParseToken(keyString, tokenType string) (string, error) {
	if keyString == "" {
		return "", utils.ErrNotFound
	}
	inputTokenType := keyString[0:len(tokenType)]
	inputTokenString := keyString[len(tokenType):]
	if inputTokenType != tokenType {
		return "", utils.ErrInvalid
	}
	return strings.TrimSpace(inputTokenString), nil
}

func JwtAuthMiddleware(tokenKey, tokenType string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		keyString := ctx.Request.Header.Get("Authorization")
		inputTokenString, err := ParseToken(keyString, tokenType)
		if err != nil {
			if err == utils.ErrNotFound {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, datamodels.DefaultResponse{
					Code:    1,
					Message: "Authorization Required!",
				})
				return
			}
			if err == utils.ErrInvalid {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, datamodels.DefaultResponse{
					Code:    2,
					Message: fmt.Sprintf("Authorization Required! [%v]", keyString[0:len(tokenType)]),
				})
				return
			}
		}
		claim := jwt.RegisteredClaims{}
		if err := ValidateToken(inputTokenString, tokenKey, &claim); err != nil {
			if errors.Is(err, jwt.ErrTokenExpired) {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, datamodels.DefaultResponse{
					Code:    3,
					Message: "Authorization Required! [TokenExpired]",
				})
			} else {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, datamodels.DefaultResponse{
					Code:    4,
					Message: fmt.Sprintf("Authorization Required! [%v]", err),
				})
			}
			return
		}
		ctx.Set("userCode", claim.Subject)
	}
}
