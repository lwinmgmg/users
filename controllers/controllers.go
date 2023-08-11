package controllers

import (
	"github.com/gin-gonic/gin"
	v1 "github.com/lwinmgmg/user/controllers/v1"
)

func DefineRoutes(app *gin.Engine) {
	v1Router := app.Group("/api/v1")

	userController := &v1.UserController{
		Router: v1Router,
	}
	userController.HandleRoutes()
}
