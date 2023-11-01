package middlewares_test

import (
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/lwinmgmg/user/middlewares"
)

func TestValidateToken(t *testing.T) {
	var defaultKey string = "letmein"
	tokenString := middlewares.GetDefaultToken("myname", defaultKey, defaultKey)
	myClaim := jwt.RegisteredClaims{}
	if err := middlewares.ValidateToken(tokenString, defaultKey, &myClaim); err != nil {
		t.Error(err)
	}
}
