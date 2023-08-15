package utils

import (
	"testing"

	"github.com/golang-jwt/jwt/v5"
)

func TestValidateToken(t *testing.T) {
	tokenString := GetDefaultToken("myname", DefaultTokenKey)
	myClaim := jwt.RegisteredClaims{}
	if err := ValidateToken(tokenString, DefaultTokenKey, &myClaim); err != nil {
		t.Error(err)
	}
}
