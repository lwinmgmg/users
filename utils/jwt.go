package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	DefaultTokenKey = "letmein"
	BearerTokenType = "Bearer"
	OtpTokenType    = "UUID"
)

type UserClaim struct {
	jwt.RegisteredClaims
}

func (userClaim *UserClaim) Valid() error {
	return nil
}

func GetToken(userCode, tokenKey string, expiresAfter time.Duration) string {
	nowTime := time.Now().UTC()
	claim := jwt.RegisteredClaims{
		Issuer:    "user",
		IssuedAt:  jwt.NewNumericDate(nowTime),
		ExpiresAt: jwt.NewNumericDate(nowTime.Add(expiresAfter)),
		Subject:   userCode,
	}
	tkn := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	output, _ := tkn.SignedString([]byte(tokenKey))
	return output
}

func GetDefaultToken(userCode, tokenKey string) string {
	return GetToken(userCode, tokenKey, 24*time.Hour)
}

func ValidateToken(tokenString, tokenKey string, claim jwt.Claims) error {
	_, err := jwt.ParseWithClaims(tokenString, claim, func(token *jwt.Token) (interface{}, error) {
		return []byte(tokenKey), nil
	})
	return err
}
