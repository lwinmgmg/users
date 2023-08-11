package v1

import (
	"github.com/lwinmgmg/user/services"
	"gorm.io/gorm"
)

var (
	DB *gorm.DB = services.PgDb
)
