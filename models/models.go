package models

import (
	"fmt"
	"time"

	"github.com/lwinmgmg/user/services"
	uuid_code "github.com/lwinmgmg/uuid_code/v1"
	"gorm.io/gorm"
)

var (
	UuidCode *uuid_code.UuidCode = uuid_code.NewDefaultUuidCode()
	DB       *gorm.DB            = services.PgDb
)

type DefaultModel struct {
	ID         uint      `gorm:"primaryKey"`
	CreateDate time.Time `gorm:"autoCreateTime:nano"`
	WriteDate  time.Time `gorm:"autoUpdateTime:nano"`
}

func init() {
	user := &User{}
	partner := &Partner{}
	DB.AutoMigrate(partner)
	DB.AutoMigrate(user)
	DB.Exec(fmt.Sprintf("CREATE SEQUENCE IF NOT EXISTS %v;", partner.GetSequence()))
	DB.Exec(fmt.Sprintf("CREATE SEQUENCE IF NOT EXISTS %v START WITH 100000;", user.GetSequence()))
}
