package controllers

import (
	"github.com/gin-gonic/gin"
	v1 "github.com/lwinmgmg/user/controllers/v1"
	"github.com/lwinmgmg/user/middlewares"
	"github.com/lwinmgmg/user/utils"
)

func DefineRoutes(app *gin.Engine) {
	v1AuthRouter := app.Group("/api/v1")

	userAuthController := &v1.UserAuthController{
		Router: v1AuthRouter,
	}
	userAuthController.HandleRoutes()
	v1Router := app.Group("/api/v1", middlewares.JwtAuthMiddleware(utils.DefaultTokenKey, utils.TokenType))
	userController := &v1.UserController{
		Router: v1Router,
	}
	userController.HandleRoutes()
}
