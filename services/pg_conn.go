package services

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	PgDsn string = "host=localhost user=lwinmgmg password=letmein dbname=gorm port=5432 sslmode=disable TimeZone=Asia/Rangoon"
)

var (
	PgDb *gorm.DB
)

func GetPgConn(dsn string) (*gorm.DB, error) {
	return gorm.Open(postgres.Open(dsn), &gorm.Config{
		PrepareStmt: true,
	})
}

func init() {
	var err error
	if PgDb == nil {
		PgDb, err = GetPgConn(PgDsn)
		if err != nil {
			panic(err)
		}
	}
}
