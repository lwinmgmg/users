package middlewares_test

import (
	"fmt"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/lwinmgmg/user/middlewares"
	"github.com/lwinmgmg/user/models"
	"github.com/lwinmgmg/user/services"
	"gorm.io/gorm"
)

func TestValidateToken(t *testing.T) {
	var defaultKey string = "letmein"
	username := "test"
	password := "password"
	user, err := models.CreateTestUser(username, password, services.PgDb)
	if err != nil {
		t.Error(err)
	}
	defer services.PgDb.Delete(&user)
	services.PgDb.Transaction(func(tx *gorm.DB) error {
		tokenString := middlewares.GetDefaultToken(user.Code, string(user.Password), defaultKey)
		myClaim := jwt.RegisteredClaims{}
		if err := middlewares.ValidateToken(tokenString, defaultKey, &myClaim); err != nil {
			t.Error(err)
		}
		return fmt.Errorf("Error to rollback")
	})

}
