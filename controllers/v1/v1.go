package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/lwinmgmg/user/controllers"
	"github.com/lwinmgmg/user/controllers/v1/user"
	"github.com/lwinmgmg/user/env"
	"github.com/lwinmgmg/user/middlewares"
	"github.com/lwinmgmg/user/services"
	"gorm.io/gorm"
)

var (
	Env          = env.GetEnv()
	DB  *gorm.DB = services.PgDb
)

func ServeController(ctrls ...controllers.HttpController) {
	for i := 0; i < len(ctrls); i++ {
		ctrls[i].HandleRoutes()
	}
}

func DefineRoutes(app *gin.Engine) {
	v1Router := app.Group("/api/v1")
	v1JwtRouter := app.Group("/api/v1", middlewares.JwtAuthMiddleware(middlewares.DefaultTokenKey, middlewares.BearerTokenType))

	ServeController(
		&user.UserAuthController{
			Router: v1Router,
			DB:     DB,
		},
		&user.UserController{
			Router: v1JwtRouter,
			DB:     DB,
		},
	)
}
