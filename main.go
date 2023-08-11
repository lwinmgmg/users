package main

import (
	"fmt"
	"sync"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/lwinmgmg/user/controllers"
	"github.com/lwinmgmg/user/models"
	"gorm.io/gorm"
)

var (
	wg sync.WaitGroup
)

func printSomething(partner *models.Partner, db *gorm.DB) {
	fmt.Println(partner.NextCode(db))
	wg.Done()
}

func main() {
	app := gin.Default()
	app.Use(cors.Default())
	controllers.DefineRoutes(app)
	app.Run("0.0.0.0:8888")
}
