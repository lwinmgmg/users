package middlewares

import (
	"github.com/lwinmgmg/user/env"
	"github.com/lwinmgmg/user/services"
)

var (
	Env  = env.GetEnv()
	PgDb = services.PgDb
)
