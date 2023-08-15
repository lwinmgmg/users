package middlewares

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/lwinmgmg/user/datamodels"
	"github.com/lwinmgmg/user/services"
	"github.com/lwinmgmg/user/utils"
	"github.com/redis/go-redis/v9"
)

func JwtAuthMiddleware(tokenKey, tokenType string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		keyString := ctx.Request.Header.Get("Authorization")
		if keyString == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, datamodels.DefaultResponse{
				Code:    1,
				Message: "Authorization Required!",
			})
			return
		}
		inputTokenType := keyString[0:len(tokenType)]
		inputTokenString := keyString[len(tokenType):]
		if inputTokenType != tokenType {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, datamodels.DefaultResponse{
				Code:    2,
				Message: fmt.Sprintf("Authorization Required! Wrong Token Type [%v]", inputTokenType),
			})
			return
		}
		inputTokenString = strings.TrimSpace(inputTokenString)
		if _, err := services.GetKey(inputTokenString); err != nil {
			claim := jwt.RegisteredClaims{}
			if tknErr := utils.ValidateToken(inputTokenString, tokenKey, &claim); tknErr != nil {
				if tknErr == jwt.ErrTokenExpired {
					ctx.AbortWithStatusJSON(http.StatusUnauthorized, datamodels.DefaultResponse{
						Code:    3,
						Message: "Authorization Required! [TokenExpired]",
					})
				} else {
					ctx.AbortWithStatusJSON(http.StatusUnauthorized, datamodels.DefaultResponse{
						Code:    4,
						Message: fmt.Sprintf("Authorization Required! [%v]", tknErr),
					})
				}
				return
			}
			if err == redis.Nil {
				services.SetKey(inputTokenString, "1", claim.ExpiresAt.Sub(time.Now()))
			}
		}
	}
}
