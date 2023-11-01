package middlewares

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/lwinmgmg/user/models"
	"github.com/lwinmgmg/user/services"
)

var (
	DefaultTokenKey = Env.Settings.JwtKey
	BearerTokenType = "Bearer"
	OtpTokenType    = "UUID"
)

type UserClaim struct {
	jwt.RegisteredClaims
}

func (userClaim *UserClaim) Valid() error {
	return nil
}

func GetToken(userCode, password, tokenKey string, expiresAfter time.Duration) string {
	nowTime := time.Now().UTC()
	claim := jwt.RegisteredClaims{
		Issuer:    "user",
		IssuedAt:  jwt.NewNumericDate(nowTime),
		ExpiresAt: jwt.NewNumericDate(nowTime.Add(expiresAfter)),
		Subject:   userCode,
	}
	tkn := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	output, _ := tkn.SignedString([]byte(password + tokenKey))
	return output
}

func GetDefaultToken(userCode, password, tokenKey string) string {
	return GetToken(userCode, password, tokenKey, 24*time.Hour)
}

func GetOrSetPassword(code string) (password string, err error) {
	keyFormat := fmt.Sprintf("%v.password", code)
	if password, err = services.GetKey(keyFormat); err != nil {
		if password, err = models.GetPasswordByUserCode(code, PgDb); err != nil {
			return password, err
		}
		if _, err = services.SetKey(keyFormat, password, 1*time.Hour); err != nil {
			return password, err
		}
	}
	return password, err
}

func ValidateToken(tokenString, tokenKey string, claim jwt.Claims) error {
	_, err := jwt.ParseWithClaims(tokenString, claim, func(token *jwt.Token) (interface{}, error) {
		sub, err := token.Claims.GetSubject()
		if err != nil {
			return nil, err
		}
		password, err := GetOrSetPassword(sub)
		if err != nil {
			return nil, err
		}
		return []byte(password + tokenKey), nil
	})
	return err
}
