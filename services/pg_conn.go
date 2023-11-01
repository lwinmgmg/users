package services

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	PgDb *gorm.DB
)

func GetPgConn(dsn string) (*gorm.DB, error) {
	return gorm.Open(postgres.Open(dsn), &gorm.Config{
		PrepareStmt: true,
	})
}
